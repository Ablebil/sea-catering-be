package fiber

import (
	"sea-catering/internal/domain/env"
	"time"

	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/wI2L/jettison"
)

func New(env *env.Env) *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
		JSONEncoder: jettison.Marshal,
		JSONDecoder: gojson.Unmarshal,
	})

	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${method} | ${path} | ${latency}\n",
	}))

	app.Use(helmet.New())

	app.Use(healthcheck.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))

	app.Use(compress.New())

	return app
}
