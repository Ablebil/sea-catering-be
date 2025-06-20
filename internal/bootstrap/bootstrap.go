package bootstrap

import (
	"fmt"

	conf "github.com/Ablebil/sea-catering-be/config"
	"github.com/Ablebil/sea-catering-be/internal/infra/fiber"
	"github.com/Ablebil/sea-catering-be/internal/infra/postgresql"
)

func Start() error {
	config, err := conf.New()
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
