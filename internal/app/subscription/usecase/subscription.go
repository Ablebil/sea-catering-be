package usecase

import (
	"fmt"
	"strings"
	"time"

	mealPlanRepository "github.com/Ablebil/sea-catering-be/internal/app/meal_plan/repository"
	subscriptionRepository "github.com/Ablebil/sea-catering-be/internal/app/subscription/repository"
	"github.com/Ablebil/sea-catering-be/internal/domain/dto"
	"github.com/Ablebil/sea-catering-be/internal/domain/entity"
	"github.com/Ablebil/sea-catering-be/internal/infra/midtrans"
	res "github.com/Ablebil/sea-catering-be/internal/infra/response"
	"github.com/google/uuid"
)

type SubscriptionUsecaseItf interface {
	CreateSubscription(userID uuid.UUID, email string, req dto.CreateSubscriptionRequest) (*dto.PaymentResponse, *res.Err)
	GetUserSubscriptions(userID uuid.UUID) ([]dto.SubscriptionResponse, *res.Err)
	PauseSubscription(userID uuid.UUID, subscriptionID uuid.UUID, req dto.PauseSubscriptionRequest) (*dto.SubscriptionResponse, *res.Err)
	CancelSubscription(userID uuid.UUID, subscriptionID uuid.UUID) (*dto.SubscriptionResponse, *res.Err)
	HandlePaymentNotification(notification map[string]interface{}) *res.Err
	UpdateExpiredSubscriptions() *res.Err
}

type SubscriptionUsecase struct {
	SubscriptionRepository subscriptionRepository.SubscriptionRepositoryItf
	MealPlanRepository     mealPlanRepository.MealPlanRepositoryItf
	midtrans               midtrans.MidtransItf
}

func NewSubscriptionUsecase(subscriptionRepository subscriptionRepository.SubscriptionRepositoryItf, mealPlanRepository mealPlanRepository.MealPlanRepositoryItf, midtrans midtrans.MidtransItf) SubscriptionUsecaseItf {
	return &SubscriptionUsecase{
		SubscriptionRepository: subscriptionRepository,
		MealPlanRepository:     mealPlanRepository,
		midtrans:               midtrans,
	}
}

func (uc *SubscriptionUsecase) CreateSubscription(userID uuid.UUID, email string, req dto.CreateSubscriptionRequest) (*dto.PaymentResponse, *res.Err) {
	mealPlan, err := uc.MealPlanRepository.GetMealPlanByID(req.MealPlanID)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetMealPlanByID)
	}

	if mealPlan == nil {
		return nil, res.ErrNotFound(res.MealPlanNotFound)
	}

	totalPrice := mealPlan.Price * float64(len(req.MealTypes)) * float64(len(req.DeliveryDays)) * 4.3

	orderID := "SUBS-" + uuid.NewString()
	now := time.Now()
	end := now.Add(30 * 24 * time.Hour)

	newSubscription := &entity.Subscription{
		UserID:       userID,
		MealPlanID:   req.MealPlanID,
		Name:         req.Name,
		PhoneNumber:  req.PhoneNumber,
		Status:       entity.StatusPending,
		MealTypes:    strings.Join(req.MealTypes, ","),
		DeliveryDays: strings.Join(req.DeliveryDays, ","),
		Allergies:    req.Allergies,
		TotalPrice:   totalPrice,
		OrderID:      &orderID,
		StartDate:    now,
		EndDate:      &end,
	}

	if err := uc.SubscriptionRepository.CreateSubscription(newSubscription); err != nil {
		return nil, res.ErrInternalServerError(res.FailedSaveSubscription)
	}

	midtransReq := &dto.MidtransRequest{
		OrderID:        orderID,
		Amount:         int64(totalPrice),
		SubscriptionID: newSubscription.ID,
		CustomerDetails: dto.MidtransCustomerDetails{
			Name:  req.Name,
			Email: email,
			Phone: req.PhoneNumber,
		},
		ItemDetails: []dto.MidtransItemDetail{{
			ID:    mealPlan.ID.String(),
			Name:  fmt.Sprintf("Subscription %s", mealPlan.Name),
			Price: int64(totalPrice),
			Qty:   1,
		}},
	}

	paymentResponse, err := uc.midtrans.CreateTransaction(midtransReq)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedCreatePaymentTransaction)
	}

	return paymentResponse, nil
}

func (uc *SubscriptionUsecase) GetUserSubscriptions(userID uuid.UUID) ([]dto.SubscriptionResponse, *res.Err) {
	subs, err := uc.SubscriptionRepository.GetAllSubscriptionByUserID(userID)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetAllSubscriptions)
	}

	result := make([]dto.SubscriptionResponse, 0, len(subs))
	for _, s := range subs {
		mealPlan, err := uc.MealPlanRepository.GetMealPlanByID(s.MealPlanID)
		if err != nil {
			return nil, res.ErrInternalServerError(res.FailedGetMealPlanByID)
		}

		mealPlanResp := dto.MealPlanResponse{
			ID:          mealPlan.ID,
			Name:        mealPlan.Name,
			Description: mealPlan.Description,
			Price:       mealPlan.Price,
			PhotoURL:    mealPlan.PhotoURL,
		}

		result = append(result, dto.SubscriptionResponse{
			ID:           s.ID,
			Name:         s.Name,
			PhoneNumber:  s.PhoneNumber,
			MealPlan:     mealPlanResp,
			MealTypes:    strings.Split(s.MealTypes, ","),
			DeliveryDays: strings.Split(s.DeliveryDays, ","),
			Allergies:    s.Allergies,
			TotalPrice:   s.TotalPrice,
			Status:       string(s.Status),
			StartDate:    s.StartDate,
			EndDate:      s.EndDate,
		})
	}

	return result, nil
}

func (uc *SubscriptionUsecase) PauseSubscription(userID uuid.UUID, subscriptionID uuid.UUID, req dto.PauseSubscriptionRequest) (*dto.SubscriptionResponse, *res.Err) {
	sub, err := uc.SubscriptionRepository.GetSubscriptionByIDAndUserID(subscriptionID, userID)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetSubscriptionByID)
	}

	if sub == nil {
		return nil, res.ErrNotFound(res.SubscriptionNotFound)
	}

	sub.PauseStartDate = &req.StartDate
	sub.PauseEndDate = &req.EndDate
	sub.Status = entity.StatusPaused

	pauseDuration := req.EndDate.Sub(req.StartDate)
	if sub.EndDate != nil {
		newEndDate := sub.EndDate.Add(pauseDuration)
		sub.EndDate = &newEndDate
	}

	if err := uc.SubscriptionRepository.UpdateSubscription(sub); err != nil {
		return nil, res.ErrInternalServerError(res.FailedPauseSubscription)
	}

	mealPlan, err := uc.MealPlanRepository.GetMealPlanByID(sub.MealPlanID)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetMealPlanByID)
	}

	mealPlanResp := dto.MealPlanResponse{
		ID:          mealPlan.ID,
		Name:        mealPlan.Name,
		Description: mealPlan.Description,
		Price:       mealPlan.Price,
		PhotoURL:    mealPlan.PhotoURL,
	}

	return &dto.SubscriptionResponse{
		ID:             sub.ID,
		Name:           sub.Name,
		PhoneNumber:    sub.PhoneNumber,
		MealPlan:       mealPlanResp,
		MealTypes:      strings.Split(sub.MealTypes, ","),
		DeliveryDays:   strings.Split(sub.DeliveryDays, ","),
		Allergies:      sub.Allergies,
		TotalPrice:     sub.TotalPrice,
		Status:         string(sub.Status),
		PauseStartDate: sub.PauseStartDate,
		PauseEndDate:   sub.PauseEndDate,
		StartDate:      sub.StartDate,
		EndDate:        sub.EndDate,
	}, nil
}

func (uc *SubscriptionUsecase) CancelSubscription(userID uuid.UUID, subscriptionID uuid.UUID) (*dto.SubscriptionResponse, *res.Err) {
	sub, err := uc.SubscriptionRepository.GetSubscriptionByIDAndUserID(subscriptionID, userID)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetSubscriptionByID)
	}

	if sub == nil {
		return nil, res.ErrNotFound(res.SubscriptionNotFound)
	}

	sub.Status = entity.StatusCancelled
	if err := uc.SubscriptionRepository.UpdateSubscription(sub); err != nil {
		return nil, res.ErrInternalServerError(res.FailedCancelSubscription)
	}

	mealPlan, err := uc.MealPlanRepository.GetMealPlanByID(sub.MealPlanID)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetMealPlanByID)
	}

	mealPlanResp := dto.MealPlanResponse{
		ID:          mealPlan.ID,
		Name:        mealPlan.Name,
		Description: mealPlan.Description,
		Price:       mealPlan.Price,
		PhotoURL:    mealPlan.PhotoURL,
	}

	return &dto.SubscriptionResponse{
		ID:           sub.ID,
		Name:         sub.Name,
		PhoneNumber:  sub.PhoneNumber,
		MealPlan:     mealPlanResp,
		MealTypes:    strings.Split(sub.MealTypes, ","),
		DeliveryDays: strings.Split(sub.DeliveryDays, ","),
		Allergies:    sub.Allergies,
		TotalPrice:   sub.TotalPrice,
		Status:       string(sub.Status),
		StartDate:    sub.StartDate,
		EndDate:      sub.EndDate,
	}, nil
}

func (uc *SubscriptionUsecase) HandlePaymentNotification(notification map[string]interface{}) *res.Err {
	orderID, ok := notification["order_id"].(string)
	if !ok {
		return res.ErrBadRequest(res.InvalidOrderID)
	}

	transactionStatus, ok := notification["transaction_status"].(string)
	if !ok {
		return res.ErrBadRequest(res.InvalidTransactionStatus)
	}

	subscription, err := uc.SubscriptionRepository.GetSubscriptionByOrderID(orderID)
	if err != nil {
		return res.ErrInternalServerError(res.FailedGetSubscriptionByID)
	}

	if subscription == nil {
		return res.ErrNotFound(res.SubscriptionNotFound)
	}

	var newStatus entity.SubscriptionStatus
	switch transactionStatus {
	case "capture", "settlement":
		newStatus = entity.StatusActive
	case "cancel", "expire", "failure":
		newStatus = entity.StatusCancelled
	case "pending":
		newStatus = entity.StatusPending
	default:
		return nil
	}

	subscription.Status = newStatus
	if err := uc.SubscriptionRepository.UpdateSubscription(subscription); err != nil {
		return res.ErrInternalServerError(res.FailedUpdateSubscription)
	}

	return nil
}

func (uc *SubscriptionUsecase) UpdateExpiredSubscriptions() *res.Err {
	expiredSubs, err := uc.SubscriptionRepository.GetExpiredActiveSubscriptions()
	if err != nil {
		return res.ErrInternalServerError(res.FailedGetExpiredSubscriptions)
	}

	for _, sub := range expiredSubs {
		sub.Status = entity.StatusFinished
		if err := uc.SubscriptionRepository.UpdateSubscription(&sub); err != nil {
			continue
		}
	}

	return nil
}
