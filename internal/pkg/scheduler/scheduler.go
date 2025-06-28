package scheduler

import (
	"log"

	subscriptionUsecase "github.com/Ablebil/sea-catering-be/internal/app/subscription/usecase"
	userUsecase "github.com/Ablebil/sea-catering-be/internal/app/user/usecase"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron                *cron.Cron
	subscriptionUsecase subscriptionUsecase.SubscriptionUsecaseItf
	userUsecase         userUsecase.UserUsecaseItf
}

func NewScheduler(subscriptionUsecase subscriptionUsecase.SubscriptionUsecaseItf, userUsecase userUsecase.UserUsecaseItf) *Scheduler {
	return &Scheduler{
		cron:                cron.New(),
		subscriptionUsecase: subscriptionUsecase,
		userUsecase:         userUsecase,
	}
}

func (s *Scheduler) Start() {
	s.cron.AddFunc("0 0 * * *", s.updateExpiredSubscriptions)
	s.cron.AddFunc("0 * * * *", s.removeUnverifiedUsers)
	s.cron.Start()
	log.Println("Scheduler started")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

func (s *Scheduler) updateExpiredSubscriptions() {
	log.Println("Updating expired subscriptions...")
	if err := s.subscriptionUsecase.UpdateExpiredSubscriptions(); err != nil {
		log.Printf("Error updating expired subscriptions: %v", err)
	}
}

func (s *Scheduler) removeUnverifiedUsers() {
	log.Println("Removing unverified users...")
	if err := s.userUsecase.RemoveUnverifiedUsers(); err != nil {
		log.Printf("Error removing unverified users: %v", err)
	}
}
