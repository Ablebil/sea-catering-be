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
	routerGroup.Post("/refresh-token", authHandler.RefreshToken)
}

// @Summary      Register User
// @Description  Create a new user account and send an OTP for verification.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload body dto.RegisterRequest true "Register Request"
// @Success      201  {object}  res.Res "Registration successful. OTP has been sent to email."
// @Failure      400  {object}  res.Err "Bad Request (e.g., validation error)"
// @Failure      409  {object}  res.Err "Conflict (e.g., email already exists)"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Router       /auth/register [post]
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

// @Summary      Verify OTP
// @Description  Verify the OTP sent to the user's email and get access/refresh tokens.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload body dto.VerifyOTPRequest true "Verify OTP Request"
// @Success      200  {object}  res.Res{payload=dto.TokenResponse} "Verification successful, tokens returned."
// @Failure      400  {object}  res.Err "Bad Request (e.g., invalid OTP, validation error)"
// @Failure      404  {object}  res.Err "User Not Found"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Router       /auth/verify-otp [post]
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

	payload := dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res.OK(ctx, payload, res.VerifyOTPSuccess)
}

// @Summary      Login User
// @Description  Authenticate a user with email and password and get access/refresh tokens.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload body dto.LoginRequest true "Login Request"
// @Success      200  {object}  res.Res{payload=dto.TokenResponse} "Login successful, tokens returned."
// @Failure      400  {object}  res.Err "Bad Request (validation error)"
// @Failure      401  {object}  res.Err "Unauthorized (invalid credentials or user not verified)"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Router       /auth/login [post]
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

	payload := dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res.OK(ctx, payload, res.LoginSuccess)
}

func (h AuthHandler) RefreshToken(ctx *fiber.Ctx) error {
	req := new(dto.RefreshTokenRequest)
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

	accessToken, refreshToken, err := h.AuthUsecase.RefreshToken(*req)
	if err != nil {
		return err
	}

	payload := dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res.OK(ctx, payload, res.RefreshTokenSuccess)
}

func (h AuthHandler) Logout(ctx *fiber.Ctx) error {
	req := new(dto.LogoutRequest)
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

	if err := h.AuthUsecase.Logout(*req); err != nil {
		return err
	}

	return res.OK(ctx, nil, res.LogoutSuccess)
}
