package repository

import (
	"errors"

	"github.com/Ablebil/sea-catering-be/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MealPlanRepositoryItf interface {
	GetAllMealPlans() ([]entity.MealPlan, error)
	GetMealPlanByID(id uuid.UUID) (*entity.MealPlan, error)
	CreateMealPlan(mealPlan *entity.MealPlan) error
}

type MealPlanRepository struct {
	db *gorm.DB
}

func NewMealPlanRepository(db *gorm.DB) MealPlanRepositoryItf {
	return &MealPlanRepository{
		db: db,
	}
}

func (r *MealPlanRepository) GetAllMealPlans() ([]entity.MealPlan, error) {
	var mealPlans []entity.MealPlan
	err := r.db.Order("created_at desc").Find(&mealPlans).Error
	return mealPlans, err
}

func (r *MealPlanRepository) GetMealPlanByID(id uuid.UUID) (*entity.MealPlan, error) {
	var mealPlan entity.MealPlan
	err := r.db.First(&mealPlan, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &mealPlan, nil
}

func (r *MealPlanRepository) CreateMealPlan(mealPlan *entity.MealPlan) error {
	return r.db.Create(mealPlan).Error
}
