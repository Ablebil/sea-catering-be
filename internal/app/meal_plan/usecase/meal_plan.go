package usecase

import (
	"fmt"
	"time"

	mealPlanRepository "github.com/Ablebil/sea-catering-be/internal/app/meal_plan/repository"
	"github.com/Ablebil/sea-catering-be/internal/domain/dto"
	"github.com/Ablebil/sea-catering-be/internal/domain/entity"
	"github.com/Ablebil/sea-catering-be/internal/infra/redis"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/google/uuid"
)

type MealPlanUsecaseItf interface {
	GetAllMealPlans() ([]dto.MealPlanResponse, *res.Err)
	GetMealPlanByID(id uuid.UUID) (*dto.MealPlanResponse, *res.Err)
	CreateMealPlan(req dto.CreateMealPlanRequest) *res.Err
}

type MealPlanUsecase struct {
	MealPlanRepository mealPlanRepository.MealPlanRepositoryItf
	redis              redis.RedisItf
}

func NewMealPlanUsecase(mealPlanRepository mealPlanRepository.MealPlanRepositoryItf, redis redis.RedisItf) MealPlanUsecaseItf {
	return &MealPlanUsecase{
		MealPlanRepository: mealPlanRepository,
		redis:              redis,
	}
}

func (uc *MealPlanUsecase) GetAllMealPlans() ([]dto.MealPlanResponse, *res.Err) {
	cacheKey := "meal_plans:all"
	var cachedMealPlans []dto.MealPlanResponse

	if err := uc.redis.GetCache(cacheKey, &cachedMealPlans); err == nil {
		return cachedMealPlans, nil
	}

	mealPlans, err := uc.MealPlanRepository.GetAllMealPlans()
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetAllMealPlans)
	}

	result := make([]dto.MealPlanResponse, 0, len(mealPlans))
	for _, m := range mealPlans {
		result = append(result, dto.MealPlanResponse{
			ID:          m.ID,
			Name:        m.Name,
			Description: m.Description,
			Price:       m.Price,
			PhotoURL:    m.PhotoURL,
		})
	}

	uc.redis.SetCache(cacheKey, result, 1*time.Hour)

	return result, nil
}

func (uc *MealPlanUsecase) GetMealPlanByID(id uuid.UUID) (*dto.MealPlanResponse, *res.Err) {
	cacheKey := fmt.Sprintf("meal_plan:%s", id.String())
	var cachedMealPlan dto.MealPlanResponse

	if err := uc.redis.GetCache(cacheKey, &cachedMealPlan); err == nil {
		return &cachedMealPlan, nil
	}

	mealPlan, err := uc.MealPlanRepository.GetMealPlanByID(id)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetMealPlanByID)
	}

	if mealPlan == nil {
		return nil, res.ErrNotFound(res.MealPlanNotFound)
	}

	result := &dto.MealPlanResponse{
		ID:          mealPlan.ID,
		Name:        mealPlan.Name,
		Description: mealPlan.Description,
		Price:       mealPlan.Price,
		PhotoURL:    mealPlan.PhotoURL,
	}

	uc.redis.SetCache(cacheKey, result, 1*time.Hour)

	return result, nil
}

func (uc *MealPlanUsecase) CreateMealPlan(req dto.CreateMealPlanRequest) *res.Err {
	newMealPlan := &entity.MealPlan{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		PhotoURL:    req.PhotoURL,
	}

	err := uc.MealPlanRepository.CreateMealPlan(newMealPlan)
	if err != nil {
		return res.ErrInternalServerError(res.FailedCreateMealPlan)
	}

	cacheKey := "meal_plans:all"
	uc.redis.DeleteCache(cacheKey)

	return nil
}
