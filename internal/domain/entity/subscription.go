package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionStatus string

const (
	StatusActive    SubscriptionStatus = "active"
	StatusPaused    SubscriptionStatus = "paused"
	StatusCancelled SubscriptionStatus = "cancelled"
	StatusPending   SubscriptionStatus = "pending"
	StatusFinished  SubscriptionStatus = "finished"
)

type Subscription struct {
	ID             uuid.UUID          `gorm:"column:id;type:char(36);primaryKey;not null"`
	UserID         uuid.UUID          `gorm:"column:user_id;type:char(36);not null"`
	User           *User              `gorm:"foreignKey:user_id;constraint:OnDelete:CASCADE"`
	MealPlanID     uuid.UUID          `gorm:"column:meal_plan_id;type:char(36);not null"`
	MealPlan       *MealPlan          `gorm:"foreignKey:meal_plan_id;constraint:OnDelete:CASCADE"`
	Name           string             `gorm:"column:name;type:varchar(255);not null"`
	PhoneNumber    string             `gorm:"column:phone_number;type:varchar(20);not null"`
	MealTypes      string             `gorm:"column:meal_types;type:text;not null"`
	DeliveryDays   string             `gorm:"column:delivery_days;type:text;not null"`
	Allergies      *string            `gorm:"column:allergies;type:text"`
	TotalPrice     float64            `gorm:"column:total_price;type:decimal(15,2);not null"`
	Status         SubscriptionStatus `gorm:"column:status;type:varchar(20);default:'pending'"`
	OrderID        *string            `gorm:"column:order_id;type:varchar(255);unique"`
	PauseStartDate *time.Time         `gorm:"column:pause_start_date;type:date"`
	PauseEndDate   *time.Time         `gorm:"column:pause_end_date;type:date"`
	StartDate      time.Time          `gorm:"column:start_date;type:date;not null"`
	EndDate        *time.Time         `gorm:"column:end_date;type:date;not null"`
	CreatedAt      *time.Time         `gorm:"column:created_at;type:timestamp;autoCreateTime"`
	UpdatedAt      *time.Time         `gorm:"column:updated_at;type:timestamp;autoUpdateTime"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewV7()
	s.ID = id
	return
}
