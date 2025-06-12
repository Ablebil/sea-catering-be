package bootstrap

import (
	"fmt"
	"sea-catering/internal/domain/env"
	"sea-catering/internal/infra/fiber"
	"sea-catering/internal/infra/postgresql"
)

func Start() error {
	config, err := env.New()
	if err != nil {
		panic(err)
	}

	db, err := postgresql.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBPort,
	), config)

	app := fiber.New(config)
	v1 := app.Group("/api/v1")

	return app.Listen(fmt.Sprintf("%s:%d", config.AppHost, config.AppPort))
}
