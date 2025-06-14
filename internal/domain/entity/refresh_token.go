package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uuid.UUID  `gorm:"column:id;type:char(36);primaryKey;not null"`
	UserID    uuid.UUID  `gorm:"column:user_id;type:char(36);not null"`
	User      *User      `gorm:"foreignKey:user_id;constraint:OnDelete:CASCADE"`
	Token     string     `gorm:"column:token;type:varchar(255);not null"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime"`
}

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewV7()
	r.ID = id
	return
}
