package middleware

import (
	"github.com/Ablebil/sea-catering-be/internal/domain/entity"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) Authorization(ctx *fiber.Ctx) error {
	role, ok := ctx.Locals("role").(entity.UserRole)
	if !ok || role != entity.RoleAdmin {
		return res.ErrForbidden(res.AdminAccessRequired)
	}

	return ctx.Next()
}
