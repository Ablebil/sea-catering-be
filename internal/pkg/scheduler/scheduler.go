package scheduler

import (
	"log"

	"github.com/Ablebil/sea-catering-be/internal/app/subscription/usecase"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron                *cron.Cron
	subscriptionUsecase usecase.SubscriptionUsecaseItf
}

func NewScheduler(subscriptionUsecase usecase.SubscriptionUsecaseItf) *Scheduler {
	return &Scheduler{
		cron:                cron.New(),
		subscriptionUsecase: subscriptionUsecase,
	}
}

func (s *Scheduler) Start() {
	s.cron.AddFunc("0 0 * * *", s.updateExpiredSubscriptions)
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
