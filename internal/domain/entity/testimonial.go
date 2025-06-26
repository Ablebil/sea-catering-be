package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Testimonial struct {
	ID        uuid.UUID  `gorm:"column:id;type:char(36);primaryKey;not null"`
	UserID    uuid.UUID  `gorm:"column:user_id;type:char(36);not null"`
	User      *User      `gorm:"foreignKey:user_id;constraint:OnDelete:CASCADE"`
	Name      string     `gorm:"column:name;type:varchar(255);not null"`
	Message   string     `gorm:"column:message;type:text;not null"`
	Rating    int        `gorm:"column:rating;type:int;not null"`
	PhotoURL  string     `gorm:"column:photo_url;type:text;not null"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime"`
}

func (t *Testimonial) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewV7()
	t.ID = id
	return
}
