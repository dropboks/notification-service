package di

import (
	"github.com/dropboks/notification-service/config/logger"
	mq "github.com/dropboks/notification-service/config/message-queue"
	"github.com/dropboks/notification-service/internal/domain/handler"
	"github.com/dropboks/notification-service/internal/domain/service"
	_mq "github.com/dropboks/notification-service/internal/infrastructure/message-queue"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	if err := container.Provide(logger.New); err != nil {
		panic("Failed to provide logger: " + err.Error())
	}
	if err := container.Provide(mq.New); err != nil {
		panic("Failed to provide message queue: " + err.Error())
	}
	if err := container.Provide(jetstream.New); err != nil {
		panic("Failed to provide jetstream: " + err.Error())
	}
	if err := container.Provide(_mq.NewNatsInfrastructure); err != nil {
		panic("Failed to provide message queue infra: " + err.Error())
	}
	if err := container.Provide(service.NewSubscriberService); err != nil {
		panic("Failed to provide subscriber service " + err.Error())
	}
	if err := container.Provide(handler.NewSubscriberHandler); err != nil {
		panic("Failed to provide subscriber handler" + err.Error())
	}
	return container
}
