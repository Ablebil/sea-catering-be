package rest

import (
	"github.com/Ablebil/sea-catering-be/internal/app/subscription/usecase"
	"github.com/Ablebil/sea-catering-be/internal/domain/dto"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/Ablebil/sea-catering-be/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	Validator           *validator.Validate
	SubscriptionUsecase usecase.SubscriptionUsecaseItf
}

func NewSubscriptionHandler(routerGroup fiber.Router, validator *validator.Validate, subscriptionUsecase usecase.SubscriptionUsecaseItf, middleware middleware.MiddlewareItf) {
	subscriptionHandler := SubscriptionHandler{
		Validator:           validator,
		SubscriptionUsecase: subscriptionUsecase,
	}

	routerGroup = routerGroup.Group("/subscriptions")
	routerGroup.Post("/", middleware.Authentication, subscriptionHandler.CreateSubscription)
	routerGroup.Get("/", middleware.Authentication, subscriptionHandler.GetUserSubscriptions)
	routerGroup.Put("/:id/pause", middleware.Authentication, subscriptionHandler.PauseSubscription)
	routerGroup.Delete("/:id", middleware.Authentication, subscriptionHandler.CancelSubscription)

	adminRouterGroup := routerGroup.Group("/admin", middleware.Authentication, middleware.Authorization)
	adminRouterGroup.Get("/stats/new", subscriptionHandler.GetNewSubscriptionsStats)
	adminRouterGroup.Get("/stats/mrr", subscriptionHandler.GetMRRStats)
	adminRouterGroup.Get("/stats/active-total", subscriptionHandler.GetTotalActiveSubscriptions)
	adminRouterGroup.Get("/stats/reactivations", subscriptionHandler.GetReactivationStats)

	routerGroup.Post("/webhook/midtrans", subscriptionHandler.HandleMidtransWebhook)
}

// @Summary      Create Subscription
// @Description  Create a new meal plan subscription with payment.
// @Tags         Subscription
// @Accept       json
// @Produce      json
// @Param        payload body dto.CreateSubscriptionRequest true "Create Subscription Request"
// @Success      201  {object}  res.Res{payload=dto.PaymentResponse} "Subscription created successfully"
// @Failure      400  {object}  res.Err "Invalid request body or validation error"
// @Failure      401  {object}  res.Err "Missing or invalid access token"
// @Failure      404  {object}  res.Err "Meal plan not found"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /subscriptions/ [post]
func (h SubscriptionHandler) CreateSubscription(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)
	email := ctx.Locals("email").(string)

	req := new(dto.CreateSubscriptionRequest)
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

	paymentResp, err := h.SubscriptionUsecase.CreateSubscription(userID, email, *req)
	if err != nil {
		return err
	}

	return res.Created(ctx, paymentResp, res.CreateSubscriptionSuccess)
}

// @Summary      Get User Subscriptions
// @Description  Retrieve all subscriptions for the authenticated user.
// @Tags         Subscription
// @Produce      json
// @Success      200  {object}  res.Res{payload=[]dto.SubscriptionResponse} "Get user subscriptions successful"
// @Failure      401  {object}  res.Err "Missing or invalid access token"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /subscriptions/ [get]
func (h SubscriptionHandler) GetUserSubscriptions(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)

	subs, err := h.SubscriptionUsecase.GetUserSubscriptions(userID)
	if err != nil {
		return err
	}

	return res.OK(ctx, subs, res.GetAllSubscriptionsSuccess)
}

// @Summary      Pause Subscription
// @Description  Temporarily pause a subscription by specifying start and end dates. The subscription duration will be extended by the pause period.
// @Tags         Subscription
// @Accept       json
// @Produce      json
// @Param        id      path  string                         true  "Subscription ID" Format(uuid)
// @Param        payload body  dto.PauseSubscriptionRequest  true  "Pause Subscription Request"
// @Success      200  {object}  res.Res{payload=dto.SubscriptionResponse} "Subscription paused successfully"
// @Failure      400  {object}  res.Err "Invalid subscription ID or request body"
// @Failure      401  {object}  res.Err "Missing or invalid access token"
// @Failure      403  {object}  res.Err "Not authorized to pause this subscription"
// @Failure      404  {object}  res.Err "Subscription not found"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /subscriptions/{id}/pause [put]
func (h SubscriptionHandler) PauseSubscription(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.ErrBadRequest(res.InvalidSubscriptionID)
	}

	req := new(dto.PauseSubscriptionRequest)
	if resErr := ctx.BodyParser(req); resErr != nil {
		return res.ErrBadRequest(res.FailedParsingRequestBody)
	}

	if resErr := h.Validator.Struct(req); resErr != nil {
		validationErrors, ok := resErr.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationErrors)
	}

	pausedSub, resErr := h.SubscriptionUsecase.PauseSubscription(userID, id, *req)
	if resErr != nil {
		return resErr
	}

	return res.OK(ctx, pausedSub, res.PauseSubscriptionSuccess)
}

// @Summary      Cancel Subscription
// @Description  Permanently cancel a subscription. This action cannot be undone.
// @Tags         Subscription
// @Produce      json
// @Param        id   path      string  true  "Subscription ID" Format(uuid)
// @Success      200  {object}  res.Res{payload=dto.SubscriptionResponse} "Subscription cancelled successfully"
// @Failure      400  {object}  res.Err "Invalid subscription ID"
// @Failure      401  {object}  res.Err "Missing or invalid access token"
// @Failure      403  {object}  res.Err "Not authorized to cancel this subscription"
// @Failure      404  {object}  res.Err "Subscription not found"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /subscriptions/{id} [delete]
func (h SubscriptionHandler) CancelSubscription(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return res.ErrBadRequest(res.InvalidSubscriptionID)
	}

	cancelledSub, resErr := h.SubscriptionUsecase.CancelSubscription(userID, id)
	if resErr != nil {
		return resErr
	}

	return res.OK(ctx, cancelledSub, res.CancelSubscriptionSuccess)
}

// @Summary      Get New Subscriptions Stats
// @Description  Get total number of new subscriptions in a date range (admin only).
// @Tags         Subscription
// @Produce      json
// @Param        start_date query string true "Start date (YYYY-MM-DD)"
// @Param        end_date   query string true "End date (YYYY-MM-DD)"
// @Success      200  {object}  res.Res{payload=object{count=int64}} "Get new subscriptions stats success"
// @Failure      400  {object}  res.Err "Invalid request params"
// @Failure      401  {object}  res.Err "Missing or invalid access token"
// @Failure      403  {object}  res.Err "Admin access required"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /subscriptions/admin/stats/new [get]
func (h SubscriptionHandler) GetNewSubscriptionsStats(ctx *fiber.Ctx) error {
	req := new(dto.GetSubscriptionStatisticRequest)
	if err := ctx.QueryParser(req); err != nil {
		return res.ErrBadRequest(res.FailedParsingRequestParams)
	}

	if err := h.Validator.Struct(req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationErrors)
	}

	count, err := h.SubscriptionUsecase.GetNewSusbcriptionsCount(*req)
	if err != nil {
		return err
	}

	return res.OK(ctx, fiber.Map{"count": count}, res.GetNewSubscriptionsStatsSuccess)
}

// @Summary      Get MRR Stats
// @Description  Get Monthly Recurring Revenue (MRR) in a date range (admin only).
// @Tags         Subscription
// @Produce      json
// @Param        start_date query string true "Start date (YYYY-MM-DD)"
// @Param        end_date   query string true "End date (YYYY-MM-DD)"
// @Success      200  {object}  res.Res{payload=object{mrr=float64}} "Get MRR stats success"
// @Failure      400  {object}  res.Err "Invalid request params"
// @Failure      401  {object}  res.Err "Missing or invalid access token"
// @Failure      403  {object}  res.Err "Admin access required"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /subscriptions/admin/stats/mrr [get]
func (h SubscriptionHandler) GetMRRStats(ctx *fiber.Ctx) error {
	req := new(dto.GetSubscriptionStatisticRequest)
	if err := ctx.QueryParser(req); err != nil {
		return res.ErrBadRequest(res.FailedParsingRequestParams)
	}

	if err := h.Validator.Struct(req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationErrors)
	}

	mrr, err := h.SubscriptionUsecase.GetMRR(*req)
	if err != nil {
		return err
	}

	return res.OK(ctx, fiber.Map{"mrr": mrr}, res.GetMRRStatsSuccess)
}

// @Summary      Get Total Active Subscriptions
// @Description  Get total number of active subscriptions (admin only).
// @Tags         Subscription
// @Produce      json
// @Success      200  {object}  res.Res{payload=object{count=int64}} "Get total active subscriptions success"
// @Failure      401  {object}  res.Err "Missing or invalid access token"
// @Failure      403  {object}  res.Err "Admin access required"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /subscriptions/admin/stats/active-total [get]
func (h SubscriptionHandler) GetTotalActiveSubscriptions(ctx *fiber.Ctx) error {
	count, err := h.SubscriptionUsecase.GetTotalActiveSubscriptions()
	if err != nil {
		return err
	}

	return res.OK(ctx, fiber.Map{"count": count}, res.GetTotalActiveSubscriptionsSuccess)
}

// @Summary      Get Reactivation Stats
// @Description  Get number of subscriptions that were reactivated in a date range (admin only).
// @Tags         Subscription
// @Produce      json
// @Param        start_date query string true "Start date (YYYY-MM-DD)"
// @Param        end_date   query string true "End date (YYYY-MM-DD)"
// @Success      200  {object}  res.Res{payload=object{count=int64}} "Get reactivation stats success"
// @Failure      400  {object}  res.Err "Invalid request params"
// @Failure      401  {object}  res.Err "Missing or invalid access token"
// @Failure      403  {object}  res.Err "Admin access required"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Security     ApiKeyAuth
// @Router       /subscriptions/admin/stats/reactivations [get]
func (h SubscriptionHandler) GetReactivationStats(ctx *fiber.Ctx) error {
	req := new(dto.GetSubscriptionStatisticRequest)
	if err := ctx.QueryParser(req); err != nil {
		return res.ErrBadRequest(res.FailedParsingRequestParams)
	}

	if err := h.Validator.Struct(req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationErrors)
	}

	count, err := h.SubscriptionUsecase.GetReactivationStats(*req)
	if err != nil {
		return err
	}

	return res.OK(ctx, fiber.Map{"count": count}, res.GetReactivationStatsSuccess)
}

// @Summary      Handle Midtrans Webhook
// @Description  Handle payment notification from Midtrans payment gateway. This endpoint is called by Midtrans to notify payment status changes.
// @Tags         Subscription
// @Accept       json
// @Produce      json
// @Param        payload body object{order_id=string,transaction_status=string,payment_type=string} true "Midtrans Notification"
// @Success      200  {object}  res.Res "Webhook processed successfully"
// @Failure      400  {object}  res.Err "Invalid notification data"
// @Failure      404  {object}  res.Err "Subscription not found"
// @Failure      500  {object}  res.Err "Internal Server Error"
// @Router       /subscriptions/webhook/midtrans [post]
func (h SubscriptionHandler) HandleMidtransWebhook(ctx *fiber.Ctx) error {
	var notification map[string]interface{}
	if err := ctx.BodyParser(&notification); err != nil {
		return res.ErrBadRequest(res.FailedParsingRequestBody)
	}

	if resErr := h.SubscriptionUsecase.HandlePaymentNotification(notification); resErr != nil {
		return resErr
	}

	return res.OK(ctx, nil, res.WebhookProcessedSuccess)
}
