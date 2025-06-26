package dto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"Str0ngP@ssword"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email" example:"john@example.com"`
	OTP   string `json:"otp" validate:"required,len=6,numeric" example:"123456"`
}

type LoginRequest struct {
	Email      string `json:"email" validate:"required,email" example:"john@example.com"`
	Password   string `json:"password" validate:"required,min=8" example:"Str0ngP@ssword"`
	RememberMe bool   `json:"remember_me" example:"true"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJI..."`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJI..."`
}

type GoogleCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
	Error string `json:"error"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJI..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJI..."`
}

type GoogleProfileResponse struct {
	ID       string `json:"google_id" example:"1234567890"`
	Email    string `json:"email" example:"john@example.com"`
	Username string `json:"username" example:"john_doe"`
	Name     string `json:"name" example:"John Doe"`
	Verified bool   `json:"verified" example:"true"`
}
