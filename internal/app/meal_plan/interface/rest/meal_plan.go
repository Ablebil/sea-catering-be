package rest

import (
	"github.com/Ablebil/sea-catering-be/internal/app/meal_plan/usecase"
	_ "github.com/Ablebil/sea-catering-be/internal/domain/dto"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MealPlanHandler struct {
	MealPlanUsecase usecase.MealPlanUsecaseItf
}

func NewMealPlanHandler(routerGroup fiber.Router, mealPlanUsecase usecase.MealPlanUsecaseItf) {
	mealPlanHandler := MealPlanHandler{
		MealPlanUsecase: mealPlanUsecase,
	}

	routerGroup = routerGroup.Group("/meal-plans")
	routerGroup.Get("/", mealPlanHandler.GetAllMealPlans)
	routerGroup.Get("/:id", mealPlanHandler.GetMealPlanByID)
}

// @Summary      Get All Meal Plans
// @Description  Get all meal plans.
// @Tags         MealPlan
// @Produce      json
// @Success      200  {object}  res.Res{payload=[]dto.MealPlanResponse} "Get all meal plans successful"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Router       /meal-plans/ [get]
func (h MealPlanHandler) GetAllMealPlans(ctx *fiber.Ctx) error {
	mealPlans, err := h.MealPlanUsecase.GetAllMealPlans()
	if err != nil {
		return err
	}

	return res.OK(ctx, mealPlans, res.GetAllMealPlansSuccess)
}

// @Summary      Get Meal Plan By ID
// @Description  Get meal plan detail by ID.
// @Tags         MealPlan
// @Produce      json
// @Param        id   path      string  true  "Meal Plan ID"
// @Success      200  {object}  res.Res{payload=dto.MealPlanResponse} "Get meal plan by ID successful"
// @Failure      400  {object}  res.Err "Invalid meal plan ID"
// @Failure      404  {object}  res.Err "Meal plan not found"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Router       /meal-plans/{id} [get]
func (h MealPlanHandler) GetMealPlanByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.ErrBadRequest(res.InvalidMealPlanID)
	}

	mealPlan, resErr := h.MealPlanUsecase.GetMealPlanByID(id)
	if resErr != nil {
		return resErr
	}

	return res.OK(ctx, mealPlan, res.GetMealPlanByIDSuccess)
}
