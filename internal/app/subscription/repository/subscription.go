package repository

import (
	"errors"
	"time"

	"github.com/Ablebil/sea-catering-be/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepositoryItf interface {
	CreateSubscription(subscription *entity.Subscription) error
	UpdateStatus(subscription *entity.Subscription, newStatus entity.SubscriptionStatus) error
	GetAllSubscriptionByUserID(userID uuid.UUID) ([]entity.Subscription, error)
	GetSubscriptionByID(id uuid.UUID) (*entity.Subscription, error)
	GetSubscriptionByIDAndUserID(id uuid.UUID, userID uuid.UUID) (*entity.Subscription, error)
	GetSubscriptionByOrderID(orderID string) (*entity.Subscription, error)
	GetExpiredActiveSubscriptions() ([]entity.Subscription, error)
	CountNewInRange(start time.Time, end time.Time) (int64, error)
	CalculateMRRInRange(start time.Time, end time.Time) (float64, error)
	CountTotalActive() (int64, error)
	CountReactivationsImage(start time.Time, end time.Time) (int64, error)
}

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepositoryItf {
	return &SubscriptionRepository{
		db: db,
	}
}

func (r *SubscriptionRepository) CreateSubscription(subscription *entity.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *SubscriptionRepository) UpdateStatus(subscription *entity.Subscription, newStatus entity.SubscriptionStatus) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		oldStatus := subscription.Status

		statusLog := &entity.SubscriptionStatusLog{
			SubscriptionID: subscription.ID,
			OldStatus:      oldStatus,
			NewStatus:      newStatus,
		}

		if err := tx.Create(statusLog).Error; err != nil {
			return err
		}

		subscription.Status = newStatus
		if err := tx.Save(subscription).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *SubscriptionRepository) GetAllSubscriptionByUserID(userID uuid.UUID) ([]entity.Subscription, error) {
	var subscriptions []entity.Subscription
	err := r.db.Preload("MealPlan").Where("user_id = ?", userID).Order("created_at desc").Find(&subscriptions).Error
	return subscriptions, err
}

func (r *SubscriptionRepository) GetSubscriptionByID(id uuid.UUID) (*entity.Subscription, error) {
	var subscription entity.Subscription
	err := r.db.Preload("MealPlan").Where("id = ?", id).First(&subscription).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *SubscriptionRepository) GetSubscriptionByIDAndUserID(id, userID uuid.UUID) (*entity.Subscription, error) {
	var subscription entity.Subscription
	err := r.db.Preload("MealPlan").Where("id = ? AND user_id = ?", id, userID).First(&subscription).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *SubscriptionRepository) GetSubscriptionByOrderID(orderID string) (*entity.Subscription, error) {
	var subscription entity.Subscription
	err := r.db.Where("order_id = ?", orderID).First(&subscription).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *SubscriptionRepository) GetExpiredActiveSubscriptions() ([]entity.Subscription, error) {
	var subscriptions []entity.Subscription
	now := time.Now()

	err := r.db.Where("status = ? AND end_date < ?", entity.StatusActive, now).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *SubscriptionRepository) CountNewInRange(start, end time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Subscription{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Count(&count).Error
	return count, err
}

func (r *SubscriptionRepository) CalculateMRRInRange(start, end time.Time) (float64, error) {
	var total float64
	err := r.db.Model(&entity.Subscription{}).
		Where("status = ? AND created_at BETWEEN ? AND ?", entity.StatusActive, start, end).
		Select("COALESCE(SUM(total_price), 0)").
		Row().Scan(&total)
	return total, err
}

func (r *SubscriptionRepository) CountTotalActive() (int64, error) {
	var count int64
	err := r.db.Model(&entity.Subscription{}).
		Where("status = ?", entity.StatusActive).
		Count(&count).Error
	return count, err
}

func (r *SubscriptionRepository) CountReactivationsImage(start time.Time, end time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&entity.SubscriptionStatusLog{}).
		Where("new_status = ? AND (old_status = ? OR old_status = ?) AND changed_at BETWEEN ? AND ?", entity.StatusActive, entity.StatusCancelled, entity.StatusFinished, start, end).
		Count(&count).Error
	return count, err
}
