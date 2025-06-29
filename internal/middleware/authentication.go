package middleware

import (
	"fmt"
	"strings"

	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) Authentication(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return res.ErrUnauthorized(res.MissingAccessToken)
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return res.ErrUnauthorized(res.InvalidOrMissingBearerToken)
	}

	userID, name, email, role, err := m.jwt.VerifyAccessToken(parts[1])
	if err != nil {
		fmt.Println("Detailed Error:", err)
		return res.ErrUnauthorized(res.InvalidAccessToken)
	}

	ctx.Locals("userID", userID)
	ctx.Locals("name", name)
	ctx.Locals("email", email)
	ctx.Locals("role", *role)

	return ctx.Next()
}
