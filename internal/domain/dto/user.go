package dto

type UserResponse struct {
	Email string `json:"email" example:"john@example.com"`
	Name  string `json:"name" example:"John Doe"`
}
