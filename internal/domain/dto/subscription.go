package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateSubscriptionRequest struct {
	Name         string    `json:"name" validate:"required,min=3,max=50" example:"John Doe"`
	PhoneNumber  string    `json:"phone_number" validate:"required,min=10" example:"081234567890"`
	MealPlanID   uuid.UUID `json:"meal_plan_id" validate:"required,uuid" example:"b3e1f8e2..."`
	MealTypes    []string  `json:"meal_types" validate:"required,min=1" example:"breakfast,lunch"`
	DeliveryDays []string  `json:"delivery_days" validate:"required,min=1" example:"monday,tuesday,wednesday"`
	Allergies    *string   `json:"allergies" example:"Peanuts, Shellfish"`
}

type PauseSubscriptionRequest struct {
	StartDate string `json:"start_date" validate:"required,datetime=2006-01-02" example:"2025-01-15"`
	EndDate   string `json:"end_date" validate:"required,datetime=2006-01-02" example:"2025-01-30"`
}

type GetSubscriptionStatisticRequest struct {
	StartDate string `query:"start_date" validate:"required,datetime=2006-01-02" example:"2025-01-15"`
	EndDate   string `query:"end_date" validate:"required,datetime=2006-01-02" example:"2025-01-30"`
}

type MidtransRequest struct {
	OrderID         string
	Amount          int64
	SubscriptionID  uuid.UUID
	ItemDetails     []MidtransItemDetail
	CustomerDetails MidtransCustomerDetails
}

type MidtransItemDetail struct {
	ID    string
	Name  string
	Price int64
	Qty   int32
}

type MidtransCustomerDetails struct {
	Name  string
	Email string
	Phone string
}

type SubscriptionResponse struct {
	ID             uuid.UUID        `json:"id" example:"b3e1f8e2..."`
	Name           string           `json:"name" example:"John Doe"`
	PhoneNumber    string           `json:"phone_number" example:"08123456789"`
	MealPlan       MealPlanResponse `json:"meal_plan"`
	MealTypes      []string         `json:"meal_types" example:"breakfast,lunch"`
	DeliveryDays   []string         `json:"delivery_days" example:"monday,tuesday,wednesday"`
	Allergies      *string          `json:"allergies" example:"Peanuts, Shellfish"`
	TotalPrice     float64          `json:"total_price" example:"180000"`
	Status         string           `json:"status" example:"pending"`
	PauseStartDate *time.Time       `json:"pause_start_date" example:"2025-01-15"`
	PauseEndDate   *time.Time       `json:"pause_end_date" example:"2025-01-30"`
	StartDate      time.Time        `json:"start_date" example:"2025-01-10"`
	EndDate        *time.Time       `json:"end_date" example:"2025-02-10"`
	CreatedAt      time.Time        `json:"created_at" example:"2025-01-10"`
}

type PaymentResponse struct {
	Token       string `json:"token" example:"66e4fa55..."`
	RedirectURL string `json:"redirect_url" example:"https://app.sandbox.midtrans.com/snap/v3/redirection/66e4fa55..."`
}
