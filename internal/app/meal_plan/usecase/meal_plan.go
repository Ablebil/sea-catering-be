package usecase

import (
	mealPlanRepository "github.com/Ablebil/sea-catering-be/internal/app/meal_plan/repository"
	"github.com/Ablebil/sea-catering-be/internal/domain/dto"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/google/uuid"
)

type MealPlanUsecaseItf interface {
	GetAllMealPlans() ([]dto.MealPlanResponse, *res.Err)
	GetMealPlanByID(id uuid.UUID) (*dto.MealPlanResponse, *res.Err)
}

type MealPlanUsecase struct {
	MealPlanRepository mealPlanRepository.MealPlanRepositoryItf
}

func NewMealPlanUsecase(mealPlanRepository mealPlanRepository.MealPlanRepositoryItf) MealPlanUsecaseItf {
	return &MealPlanUsecase{
		MealPlanRepository: mealPlanRepository,
	}
}

func (uc *MealPlanUsecase) GetAllMealPlans() ([]dto.MealPlanResponse, *res.Err) {
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

	return result, nil
}

func (uc *MealPlanUsecase) GetMealPlanByID(id uuid.UUID) (*dto.MealPlanResponse, *res.Err) {
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

	return result, nil
}
