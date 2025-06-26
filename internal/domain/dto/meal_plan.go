package dto

import "github.com/google/uuid"

type MealPlanResponse struct {
	ID          uuid.UUID `json:"id" example:"b3e1f8e2..."`
	Name        string    `json:"name" example:"John Doe"`
	Description string    `json:"description" example:"A healthy meal plan"`
	Price       float64   `json:"price" example:"30000"`
	PhotoURL    string    `json:"photo_url" example:"https://..."`
}
