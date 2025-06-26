package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MealPlan struct {
	ID          uuid.UUID  `gorm:"column:id;type:char(36);primaryKey;not null"`
	Name        string     `gorm:"column:name;type:varchar(255);not null"`
	Description string     `gorm:"column:description;type:text;not null"`
	Price       float64    `gorm:"column:price;type:decimal(10,2);not null"`
	PhotoURL    string     `gorm:"column:photo_url;type:text;not null"`
	CreatedAt   *time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime"`
}

func (m *MealPlan) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewV7()
	m.ID = id
	return
}
