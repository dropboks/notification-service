package service

import (
	mq "github.com/dropboks/notification-service/internal/infrastructure/message-queue"
	"github.com/rs/zerolog"
)

type (
	SubscriberService interface{}
	subscriberService struct {
		logger       zerolog.Logger
		natsInstance mq.Nats
	}
)

func NewSubscriberService(natsIntance mq.Nats, logger zerolog.Logger) SubscriberService {
	return &subscriberService{
		logger:       logger,
		natsInstance: natsIntance,
	}
}
