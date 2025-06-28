package rest

import (
	userUsecase "github.com/Ablebil/sea-catering-be/internal/app/user/usecase"
	_ "github.com/Ablebil/sea-catering-be/internal/domain/dto"
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

// @Summary      Get User Profile
// @Description  Get the authenticated user's profile.
// @Tags         User
// @Produce      json
// @Success      200  {object}  res.Res{payload=dto.UserResponse} "Get profile successful"
// @Failure      401  {object}  res.Err "Unauthorized"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /users/profile [get]
func (h UserHandler) GetProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)

	profile, err := h.UserUsecase.GetProfile(userID)
	if err != nil {
		return err
	}

	return res.OK(ctx, profile, res.GetProfileSuccess)
}
