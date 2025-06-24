package handler

import (
	"fmt"

	"github.com/dropboks/notification-service/internal/domain/service"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
)

type (
	SubscriberHandler interface {
		EmailHandler(msg jetstream.Msg)
	}
	subscriberHandler struct {
		subsService service.SubscriberService
		logger      zerolog.Logger
	}
)

func NewSubscriberHandler(svc service.SubscriberService, logger zerolog.Logger) SubscriberHandler {
	return &subscriberHandler{
		subsService: svc,
		logger:      logger,
	}
}

func (s *subscriberHandler) EmailHandler(msg jetstream.Msg) {
	fmt.Println(msg)
}
