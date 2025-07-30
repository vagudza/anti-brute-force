package test

import (
	"os"
	"testing"

	"github.com/vagudza/anti-brute-force/test/suitex"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	logger := initLogger()

	err := suitex.InitSuiteFactory()
	if err != nil {
		logger.Fatal("failed to initialize test suite factory", zap.Error(err))
	}

	exitVal := m.Run()
	suitex.Cleanup()
	os.Exit(exitVal)
}

func initLogger() *zap.Logger {
	return zap.Must(zap.NewDevelopment())
}
