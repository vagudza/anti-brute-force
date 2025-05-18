package test

import (
	"log"
	"os"
	"testing"

	"github.com/vagudza/anti-brute-force/test/suitex"
)

func TestMain(m *testing.M) {
	err := suitex.InitSuiteFactory()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	exitVal := m.Run()

	suitex.Cleanup()

	os.Exit(exitVal)
}
