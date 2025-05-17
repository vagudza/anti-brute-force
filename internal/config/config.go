package config

import "github.com/ilyakaznacheev/cleanenv"

type AppConfig struct {
	Env      string     `yaml:"env" env:"ENV"`
	Limiters Limiters   `yaml:"limiters" env-prefix:"LIMITERS_"`
	Postgres PGConfig   `yaml:"postgres" env-prefix:"POSTGRES_"`
	Grpc     GrpcConfig `yaml:"grpc"`
}

func New() (*AppConfig, error) {
	const defaultConfigFile = "config/app/config.local.yaml"
	var cfg AppConfig

	// ReadConfig reads config from file and override it if environment variables are set
	// File used for local development, environment variables are used in environments
	err := cleanenv.ReadConfig(defaultConfigFile, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
