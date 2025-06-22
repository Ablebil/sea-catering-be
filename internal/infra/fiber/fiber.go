package fiber

import (
	"time"

	conf "github.com/Ablebil/sea-catering-be/config"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/wI2L/jettison"
)

func New(conf *conf.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
		JSONEncoder: jettison.Marshal,
		JSONDecoder: gojson.Unmarshal,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			if customErr, ok := err.(*res.Err); ok {
				return ctx.Status(customErr.Code).JSON(customErr)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		},
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
