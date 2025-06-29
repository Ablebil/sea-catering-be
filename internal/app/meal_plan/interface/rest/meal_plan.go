package rest

import (
	"github.com/Ablebil/sea-catering-be/internal/app/meal_plan/usecase"
	"github.com/Ablebil/sea-catering-be/internal/domain/dto"
	_ "github.com/Ablebil/sea-catering-be/internal/domain/dto"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MealPlanHandler struct {
	Validator       *validator.Validate
	MealPlanUsecase usecase.MealPlanUsecaseItf
}

func NewMealPlanHandler(routerGroup fiber.Router, validator *validator.Validate, mealPlanUsecase usecase.MealPlanUsecaseItf) {
	mealPlanHandler := MealPlanHandler{
		Validator:       validator,
		MealPlanUsecase: mealPlanUsecase,
	}

	routerGroup = routerGroup.Group("/meal-plans")
	routerGroup.Get("/", mealPlanHandler.GetAllMealPlans)
	routerGroup.Get("/:id", mealPlanHandler.GetMealPlanByID)
	routerGroup.Post("/", mealPlanHandler.CreateMealPlan)
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

// @Summary      Create Meal Plan
// @Description  Create a new meal plan.
// @Tags         MealPlan
// @Accept       json
// @Produce      json
// @Param        payload body dto.CreateMealPlanRequest true "Create Meal Plan Request"
// @Success      201  {object}  res.Res "Meal plan created successfully"
// @Failure      400  {object}  res.Err "Invalid request body or validation error"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Router       /meal-plans/ [post]
func (h MealPlanHandler) CreateMealPlan(ctx *fiber.Ctx) error {
	req := new(dto.CreateMealPlanRequest)
	if err := ctx.BodyParser(req); err != nil {
		return res.ErrBadRequest(res.FailedParsingRequestBody)
	}

	if err := h.Validator.Struct(req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationErrors)
	}

	if err := h.MealPlanUsecase.CreateMealPlan(*req); err != nil {
		return err
	}

	return res.Created(ctx, nil, res.CreateMealPlanSuccess)
}
