package dto

import "github.com/google/uuid"

type CreateMealPlanRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=255" example:"Diet Plan"`
	Description string  `json:"description" validate:"required" example:"A healthy meal plan"`
	Price       float64 `json:"price" validate:"required,gt=0" example:"30000"`
	PhotoURL    string  `json:"photo_url" validate:"required,url" example:"https://..."`
}

type MealPlanResponse struct {
	ID          uuid.UUID `json:"id" example:"b3e1f8e2..."`
	Name        string    `json:"name" example:"Diet Plan"`
	Description string    `json:"description" example:"A healthy meal plan"`
	Price       float64   `json:"price" example:"30000"`
	PhotoURL    string    `json:"photo_url" example:"https://..."`
}
