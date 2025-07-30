package config

type GrpcConfig struct {
	Port string `yaml:"port" env:"GRPC_PORT"`
}
