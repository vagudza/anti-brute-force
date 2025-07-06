package suitex

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	pb "github.com/vagudza/anti-brute-force/api/proto"
	"github.com/vagudza/anti-brute-force/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var factory *suiteFactory

type suiteFactory struct {
	cfg                  *config.AppConfig
	cc                   *grpc.ClientConn
	antiBruteforceClient pb.AntiBruteforceClient
}

func (f *suiteFactory) newSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()

	const defaultTimeout = 10 * time.Second
	ctx, cancelCtx := context.WithTimeout(context.Background(), defaultTimeout)

	suite := &Suite{
		Cfg:                  f.cfg,
		AntiBruteforceClient: f.antiBruteforceClient,
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
	// Set default config path for local run from IDE
	// P.S. command "task integration-test" automatically sets CONFIG_PATH
	if os.Getenv("CONFIG_PATH") == "" {
		err := os.Setenv("CONFIG_PATH", "../config/app/config.local.yaml")
		if err != nil {
			return err
		}
	}

	cfg, err := config.New()
	if err != nil {
		return err
	}

	cc, err := grpc.NewClient(
		fmt.Sprintf(":%s", cfg.Grpc.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("can't connect to grpc server: %w", err)
	}

	factory = &suiteFactory{
		cfg:                  cfg,
		cc:                   cc,
		antiBruteforceClient: pb.NewAntiBruteforceClient(cc),
	}

	return nil
}

func Cleanup() {
	if factory != nil && factory.cc != nil {
		err := factory.cc.Close()
		if err != nil {
			log.Println(err)
		}
	}
}
