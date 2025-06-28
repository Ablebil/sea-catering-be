package dto

import (
	"github.com/google/uuid"
)

type CreateTestimonialRequest struct {
	Name    string `form:"name" validate:"required,min=3,max=50" example:"John Doe"`
	Message string `form:"message" validate:"required,min=10,max=500" example:"The food was delicious"`
	Rating  int    `form:"rating" validate:"required,min=1,max=5" example:"5"`
}

type TestimonialResponse struct {
	ID       uuid.UUID `json:"id" example:"b3e1f8e2..."`
	Name     string    `json:"name" example:"John Doe"`
	Message  string    `json:"message" example:"The food was delicious"`
	Rating   int       `json:"rating" example:"5"`
	PhotoURL string    `json:"photo_url" example:"https://..."`
}
