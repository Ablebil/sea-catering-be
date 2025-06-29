package limiter

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func Global() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	})
}

func Register() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
	})
}

func Login() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
	})
}

func VerifyOTP() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			type body struct {
				Email string `json:"email"`
			}
			var b body
			c.BodyParser(&b)
			return c.IP() + ":" + b.Email
		},
	})
}

func Testimonial() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        3,
		Expiration: 1 * time.Minute,
	})
}

func Subscription() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        3,
		Expiration: 1 * time.Minute,
	})
}
