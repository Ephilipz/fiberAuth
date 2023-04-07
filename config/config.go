package config

import (
	"github.com/caarlos0/env/v7"
)

type DB struct {
	HOST              string `env:"HOST"`
	USER              string `env:"USER"`
	PASS              string `env:"PASSWORD"`
	PORT              string `env:"PORT"`
	NAME              string `env:"NAME"`
	ENABLE_MIGRATIONS bool   `env:"ENABLE_MIGRATIONS" envDefault:"true"`
	SSLMODE           string `env:"SSLMODE" envDefault:"allow"`
}

type JWT struct {
	RSA string `env:"RSA"`
}

type App struct {
	PORT uint `env:"PORT" envDefault:"3000"`
}

type Config struct {
	Database DB  `envPrefix:"DB_"`
	Jwt      JWT `envPrefix:"JWT_"`
	App      App
}

func Get() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
