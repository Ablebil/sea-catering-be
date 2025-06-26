package rest

import (
	"github.com/Ablebil/sea-catering-be/internal/app/meal_plan/usecase"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MealPlanHandler struct {
	MealPlanUsecase usecase.MealPlanUseCaseItf
}

func NewMealPlanHandler(routerGroup fiber.Router, mealPlanUsecase usecase.MealPlanUseCaseItf) {
	mealPlanHandler := MealPlanHandler{
		MealPlanUsecase: mealPlanUsecase,
	}

	routerGroup = routerGroup.Group("/meal-plans")
	routerGroup.Get("/", mealPlanHandler.GetAllMealPlans)
	routerGroup.Get("/:id", mealPlanHandler.GetMealPlanByID)
}

func (h MealPlanHandler) GetAllMealPlans(ctx *fiber.Ctx) error {
	mealPlans, err := h.MealPlanUsecase.GetAllMealPlans()
	if err != nil {
		return err
	}

	return res.OK(ctx, mealPlans, res.GetAllMealPlansSuccess)
}

func (h MealPlanHandler) GetMealPlanByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.ErrBadRequest(res.InvalidMealPlanID)
	}

	mealPlan, err := h.MealPlanUsecase.GetMealPlanByID(id)
	if err != nil {
		return err
	}

	return res.OK(ctx, mealPlan, res.GetMealPlanByIDSuccess)
}
