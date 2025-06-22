package rest

import (
	"github.com/Ablebil/sea-catering-be/internal/app/auth/usecase"
	"github.com/Ablebil/sea-catering-be/internal/domain/dto"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Validator   *validator.Validate
	AuthUsecase usecase.AuthUsecaseItf
}

func NewAuthHandler(routerGroup fiber.Router, validator *validator.Validate, authUsecase usecase.AuthUsecaseItf) {
	authHandler := AuthHandler{
		Validator:   validator,
		AuthUsecase: authUsecase,
	}

	routerGroup = routerGroup.Group("/auth")
	routerGroup.Post("/register", authHandler.Register)
	routerGroup.Post("/verify-otp", authHandler.VerifyOTP)
	routerGroup.Post("/login", authHandler.Login)
}

func (h AuthHandler) Register(ctx *fiber.Ctx) error {
	req := new(dto.RegisterRequest)
	if err := ctx.BodyParser(req); err != nil {
		return res.ErrBadRequest(res.FailedParsingRequestBody)
	}

	if err := h.Validator.Struct(req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationErrors)
	}

	if err := h.AuthUsecase.Register(*req); err != nil {
		return err
	}

	return res.Created(ctx, nil, res.RegisterSuccess)
}

func (h AuthHandler) VerifyOTP(ctx *fiber.Ctx) error {
	req := new(dto.VerifyOTPRequest)
	if err := ctx.BodyParser(req); err != nil {
		return res.ErrBadRequest(res.FailedParsingRequestBody)
	}

	if err := h.Validator.Struct(req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationErrors)
	}

	accessToken, refreshToken, err := h.AuthUsecase.VerifyOTP(*req)
	if err != nil {
		return err
	}

	return res.OK(ctx, fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, res.VerifyOTPSuccess)
}

func (h AuthHandler) Login(ctx *fiber.Ctx) error {
	req := new(dto.LoginRequest)
	if err := ctx.BodyParser(req); err != nil {
		return res.ErrBadRequest(res.FailedParsingRequestBody)
	}

	if err := h.Validator.Struct(req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationErrors)
	}

	accessToken, refreshToken, err := h.AuthUsecase.Login(*req)
	if err != nil {
		return err
	}

	return res.OK(ctx, fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, res.LoginSuccess)
}
