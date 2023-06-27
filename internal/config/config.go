// Package config provides functions for parsing configuration files for a Go web service.
package config

// Config struct
type Config struct {
	HTTPAddr string `default:"1323"`
	DBAddr   string `default:"postgres://eugen:ur2qly1ini@localhost:5432/eugen"`
}
