// Package config provides functions for parsing configuration files for a Go web service.
package config

import (
	"github.com/caarlos0/env/v9"
)

// Config struct
type Config struct {
	HTTPAddr    string `env:"HTTP_PORT" envDefault:":8080"`
	PgxDBAddr   string `env:"PGXCONN" envDefault:"postgres://eugen:ur2qly1ini@localhost:5432/eugen"`
	MongoDBAddr string `env:"MONGODBCONN" envDefault:"mongodb://localhost:27017"`
	RedisDBAddr string `env:"REDISCONN" envDefault:"redis://:@localhost:6379/1"`
	SigningKey  string `env:"SIGNING_KEY" envDefault:"gyewgb2rf8r2b8437frb23f2er243"`
}

// NewConfig creates a new Config instance
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
