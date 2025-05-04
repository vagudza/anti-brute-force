package config

import "github.com/ilyakaznacheev/cleanenv"

type AppConfig struct {
	Postgres PGConfig `yaml:"postgres"`
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
