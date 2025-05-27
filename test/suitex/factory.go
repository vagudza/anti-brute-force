package suitex

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	pb "github.com/vagudza/anti-brute-force/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var factory *suiteFactory

type suiteFactory struct {
	antiBruteforceClient pb.AntiBruteforceClient
	cc                   *grpc.ClientConn
}

func (f *suiteFactory) newSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	const defaultTimeout = 10 * time.Second
	ctx, cancelCtx := context.WithTimeout(context.Background(), defaultTimeout)

	suite := &Suite{
		T:                    t,
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
	cfg, err := NewTestConfig()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	cc, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.GRPC.Host, cfg.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("can't connect to grpc server: %w", err)
	}

	factory = &suiteFactory{
		antiBruteforceClient: pb.NewAntiBruteforceClient(cc),
		cc:                   cc,
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
