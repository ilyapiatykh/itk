package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	DB struct {
		User     string `env:"USER,required"`
		Password string `env:"PASSWORD,required"`
		DBName   string `env:"DBNAME,required"`
	}

	Server struct {
		ReadTimeout time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
		Port        string        `env:"PORT" envDefault:":80"`
	}

	Config struct {
		DB     `envPrefix:"POSTGRES_"`
		Server `envPrefix:"SERVER_"`
	}
)

func NewCfg() (*Config, error) {
	var cfg Config
	return &cfg, env.Parse(&cfg)
}
