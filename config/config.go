package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	DB struct {
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
		DBName   string `env:"POSTGRES_DB"`
	}

	Server struct {
		ReadTimeout time.Duration `env:"SERVER_READ_TIMEOUT"`
		Port        string        `env:"PORT"`
	}

	Config struct {
		DB
		Server
	}
)

func NewCfg() (*Config, error) {
	var cfg Config
	return &cfg, env.Parse(&cfg)
}
