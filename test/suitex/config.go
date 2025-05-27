package suitex

import "github.com/ilyakaznacheev/cleanenv"

// TestConfig holds configuration for integration tests
type TestConfig struct {
	GRPC struct {
		Host string `env:"GRPC_HOST" env-default:"localhost"`
		Port string `env:"GRPC_PORT" env-default:"13013"`
	}
	Postgres struct {
		Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
		Port     string `env:"POSTGRES_PORT" env-default:"5432"`
		User     string `env:"POSTGRES_USER" env-default:"postgres"`
		Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
		Database string `env:"POSTGRES_DB" env-default:"anti-brute-force"`
	}
}

// NewTestConfig creates a new test configuration from environment
func NewTestConfig() (*TestConfig, error) {
	cfg := &TestConfig{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
