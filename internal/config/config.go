package config

type Config struct {
	HttpAddr string `default:"1323"`
	DBAddr   string `default:"postgres://eugen:ur2qly1ini@localhost:5432/eugen"`
}
