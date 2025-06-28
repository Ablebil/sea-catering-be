package rest

import (
	userUsecase "github.com/Ablebil/sea-catering-be/internal/app/user/usecase"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/Ablebil/sea-catering-be/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	UserUsecase userUsecase.UserUsecaseItf
}

func NewUserHandler(routerGroup fiber.Router, userUsecase userUsecase.UserUsecaseItf, middleware middleware.MiddlewareItf) {
	userHandler := &UserHandler{
		UserUsecase: userUsecase,
	}

	routerGroup = routerGroup.Group("/users")
	routerGroup.Get("/profile", middleware.Authentication, userHandler.GetProfile)
}

func (h UserHandler) GetProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)

	profile, err := h.UserUsecase.GetProfile(userID)
	if err != nil {
		return err
	}

	return res.OK(ctx, profile, res.GetProfileSuccess)
}
