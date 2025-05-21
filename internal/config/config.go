package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	Env      string     `yaml:"env" env:"ENV"`
	Limiters Limiters   `yaml:"limiters" env-prefix:"LIMITERS_"`
	Postgres PGConfig   `yaml:"postgres" env-prefix:"POSTGRES_"`
	Grpc     GrpcConfig `yaml:"grpc"`
}

func New() (*AppConfig, error) {
	var cfg AppConfig

	// Priority order:
	// 1. Environment variables
	// 2. Config file (if specified via CONFIG_PATH)
	// 3. Default values

	// First try to load from config file if path is specified
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		// Always read environment variables as they have highest priority
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, fmt.Errorf("error reading env variables: %w", err)
		}
	}

	// Validate config
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	return &cfg, nil
}

func validateConfig(cfg *AppConfig) error {
	if cfg.Env == "" {
		return fmt.Errorf("env must be set")
	}

	// Validate Limiters
	if cfg.Limiters.Login.MaxAttemptsPerMinute <= 0 {
		return fmt.Errorf("login max attempts per minute must be positive")
	}
	if cfg.Limiters.Password.MaxAttemptsPerMinute <= 0 {
		return fmt.Errorf("password max attempts per minute must be positive")
	}
	if cfg.Limiters.IP.MaxAttemptsPerMinute <= 0 {
		return fmt.Errorf("ip max attempts per minute must be positive")
	}

	// Validate cleanup intervals and TTL
	if cfg.Limiters.Login.CleanupInterval <= 0 {
		return fmt.Errorf("login cleanup interval must be positive")
	}
	if cfg.Limiters.Password.CleanupInterval <= 0 {
		return fmt.Errorf("password cleanup interval must be positive")
	}
	if cfg.Limiters.IP.CleanupInterval <= 0 {
		return fmt.Errorf("ip cleanup interval must be positive")
	}

	if cfg.Limiters.Login.TTL <= 0 {
		return fmt.Errorf("login TTL must be positive")
	}
	if cfg.Limiters.Password.TTL <= 0 {
		return fmt.Errorf("password TTL must be positive")
	}
	if cfg.Limiters.IP.TTL <= 0 {
		return fmt.Errorf("ip TTL must be positive")
	}

	// Validate Postgres connection
	if cfg.Postgres.Host == "" {
		return fmt.Errorf("postgres host must be set")
	}

	pgPort, err := strconv.Atoi(cfg.Postgres.Port)
	if err != nil || pgPort <= 0 {
		return fmt.Errorf("postgres port must be a positive number")
	}

	if cfg.Postgres.Database == "" {
		return fmt.Errorf("postgres database must be set")
	}
	if cfg.Postgres.Username == "" {
		return fmt.Errorf("postgres username must be set")
	}
	if cfg.Postgres.Password == "" {
		return fmt.Errorf("postgres password must be set")
	}

	// Validate GRPC
	grpcPort, err := strconv.Atoi(cfg.Grpc.Port)
	if err != nil || grpcPort <= 0 {
		return fmt.Errorf("grpc port must be a positive number")
	}

	return nil
}
