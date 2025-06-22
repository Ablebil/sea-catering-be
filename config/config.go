package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string `env:"APP_ENV"`
	AppHost string `env:"APP_HOST"`
	AppPort int    `env:"APP_PORT"`
	AppUrl  string `env:"APP_URL"`

	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`

	AccessSecret  string `env:"ACCESS_SECRET"`
	RefreshSecret string `env:"REFRESH_SECRET"`

	EmailUser     string `env:"EMAIL_USER"`
	EmailPassword string `env:"EMAIL_PASSWORD"`

	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     int    `env:"REDIS_PORT"`
	RedisPassword string `env:"REDIS_PASSWORD"`
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
