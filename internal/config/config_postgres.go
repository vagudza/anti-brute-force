package config

type PGConfig struct {
	Database string `yaml:"database" env:"DATABASE"`
	Username string `yaml:"username" env:"USERNAME"`
	Password string `yaml:"password" env:"PASSWORD"`
	Host     string `yaml:"host" env:"HOST"`
	Port     string `yaml:"port" env:"PORT"`
	SSLMode  string `yaml:"sslmode" env:"SSLMODE"`
}
