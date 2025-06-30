package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionStatusLog struct {
	ID             uuid.UUID          `gorm:"column:id;type:char(36);primaryKey;not null"`
	SubscriptionID uuid.UUID          `gorm:"column:subscription_id;type:char(36);not null"`
	Subscription   *Subscription      `gorm:"foreignKey:subscription_id;constraint:OnDelete:CASCADE"`
	OldStatus      SubscriptionStatus `gorm:"column:old_status;type:varchar(20)"`
	NewStatus      SubscriptionStatus `gorm:"column:new_status;type:varchar(20);not null"`
	ChangedAt      *time.Time         `gorm:"column:changed_at;type:timestamp;autoCreateTime"`
}

func (sl *SubscriptionStatusLog) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewV7()
	sl.ID = id
	return
}
