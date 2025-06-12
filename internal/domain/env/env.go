package env

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Env struct {
	AppEnv  string `env:"APP_ENV"`
	AppHost string `env:"APP_HOST"`
	AppPort string `env:"APP_PORT"`
	AppUrl  string `env:"APP_URLT"`

	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
}

func New() (*Env, error) {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	_env := new(Env)
	if err := env.Parse(_env); err != nil {
		return nil, err
	}

	return _env, nil
}
