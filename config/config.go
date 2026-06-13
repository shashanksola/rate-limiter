package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server ServerConfig
	Redis  RedisConfig
}

type ServerConfig struct {
	Port         int           `envconfig:"SERVER_PORT" default:"8080"`
	ReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"15s"`
	WriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"15s"`
	IdleTimeout  time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"60s"`
}

type RedisConfig struct {
	Addr     string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Load server config
	if err := envconfig.Process("", &cfg.Server); err != nil {
		return nil, fmt.Errorf("loading server config: %w", err)
	}

	// Load Redis config
	if err := envconfig.Process("", &cfg.Redis); err != nil {
		return nil, fmt.Errorf("loading redis config: %w", err)
	}

	return cfg, nil
}
