package bootstrap

import (
	"fmt"

	conf "github.com/Ablebil/sea-catering-be/config"
	"github.com/Ablebil/sea-catering-be/internal/infra/email"
	"github.com/Ablebil/sea-catering-be/internal/infra/fiber"
	"github.com/Ablebil/sea-catering-be/internal/infra/jwt"
	"github.com/Ablebil/sea-catering-be/internal/infra/oauth"
	"github.com/Ablebil/sea-catering-be/internal/infra/postgresql"
	"github.com/Ablebil/sea-catering-be/internal/infra/redis"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/swagger"

	AuthHandler "github.com/Ablebil/sea-catering-be/internal/app/auth/interface/rest"
	AuthUsecase "github.com/Ablebil/sea-catering-be/internal/app/auth/usecase"
	UserRepository "github.com/Ablebil/sea-catering-be/internal/app/user/repository"
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

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	if err := postgresql.Migrate(db); err != nil {
		return err
	}

	validator := validator.New()
	jwt := jwt.NewJWT(config)
	email := email.NewEmail(config)
	redis := redis.NewRedis(config)
	oauth := oauth.NewOAuth(config)

	app := fiber.New(config)
	v1 := app.Group("/api/v1")

	// Auth Domain
	userRepository := UserRepository.NewUserRepository(db)
	authUsecase := AuthUsecase.NewAuthUsecase(userRepository, db, config, jwt, email, redis, oauth)
	AuthHandler.NewAuthHandler(v1, validator, authUsecase, config)

	// Swagger Documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	return app.Listen(fmt.Sprintf("%s:%d", config.AppHost, config.AppPort))
}
