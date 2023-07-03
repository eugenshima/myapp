// Package config provides functions for parsing configuration files for a Go web service.
package config

// Config struct
type Config struct {
	HTTPAddr    string `env:"HTTP_PORT" envDefault:":8080"`
	PgxDBAddr   string `env:"PGXCONN" envDefault:"postgres://eugen:ur2qly1ini@localhost:5432/eugen"`
	MongoDBADDR string `env:"MONGODBCONN" envDefault:"mongodb://localhost:27017"`
}
