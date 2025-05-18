package suitex

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	pb "github.com/vagudza/anti-brute-force/api/proto"
	"github.com/vagudza/anti-brute-force/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var factory *suiteFactory

type suiteFactory struct {
	antiBrutforceClient pb.AntiBruteforceClient
	cc                  *grpc.ClientConn
}

func (f *suiteFactory) newSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	const defaultTimeout = 10 * time.Second
	ctx, cancelCtx := context.WithTimeout(context.Background(), defaultTimeout)

	suite := &Suite{
		T:                   t,
		AntiBrutforceClient: f.antiBrutforceClient,
	}

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	return ctx, suite
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	return factory.newSuite(t)
}

func InitSuiteFactory() error {
	const defaultConfigFile = "../config/app/config.local.yaml"
	var cfg config.AppConfig

	err := cleanenv.ReadConfig(defaultConfigFile, &cfg)
	if err != nil {
		return err
	}

	cc, err := grpc.NewClient(
		fmt.Sprintf("localhost:%s", cfg.Grpc.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("can't connect to grpc server: %w", err)
	}

	factory = &suiteFactory{
		antiBrutforceClient: pb.NewAntiBruteforceClient(cc),
		cc:                  cc,
	}

	return nil
}

func Cleanup() {
	err := factory.cc.Close()
	if err != nil {
		log.Println(err)
	}
}
