package config

import "time"

type Limiters struct {
	Login    LimiterConfig `yaml:"login" env-prefix:"LOGIN_"`
	Password LimiterConfig `yaml:"password" env-prefix:"PASSWORD_"`
	IP       LimiterConfig `yaml:"ip" env-prefix:"IP_"`
}

type LimiterConfig struct {
	MaxAttemptsPerMinute int           `yaml:"maxAttemptsPerMinute" env:"MAX_ATTEMPTS_PER_MINUTE"`
	CleanupInterval      time.Duration `yaml:"cleanupInterval" env:"CLEANUP_INTERVAL"`
	TTL                  time.Duration `yaml:"ttl" env:"TTL"`
}
