package postgresql

import (
	"github.com/Ablebil/sea-catering-be/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.RefreshToken{},
		&entity.Testimonial{},
		&entity.MealPlan{},
		&entity.Subscription{},
		&entity.SubscriptionStatusLog{},
	)
}
