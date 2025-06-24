package server

import (
	"context"
	"log"

	"github.com/dropboks/notification-service/internal/domain/handler"
	mq "github.com/dropboks/notification-service/internal/infrastructure/message-queue"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.uber.org/dig"
)

type Subscriber struct {
	Container       *dig.Container
	ConnectionReady chan bool
}

func (s *Subscriber) Run(ctx context.Context) {
	err := s.Container.Invoke(func(
		logger zerolog.Logger,
		sh handler.SubscriberHandler,
		js jetstream.JetStream,
		mq mq.Nats,
		_mq *nats.Conn,
	) {
		err := mq.CreateOrUpdateNewStream(ctx, &jetstream.StreamConfig{
			Name:        viper.GetString("jetstream.stream.name"),
			Description: viper.GetString("jetstream.stream.description"),
			Subjects:    []string{viper.GetString("jetstream.subject.global")},
			MaxBytes:    10 * 1024 * 1024,
			Storage:     jetstream.FileStorage,
		})
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create or update notification stream")
		}

		stream, err := js.Stream(ctx, viper.GetString("jetstream.stream.name"))
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to get stream")
		}

		// consumer for email
		emailCons, err := mq.CreateOrUpdateNewConsumer(ctx, stream, &jetstream.ConsumerConfig{
			Name:          viper.GetString("jetstream.consumer.email"),
			Durable:       viper.GetString("jetstream.consumer.email"),
			FilterSubject: viper.GetString("jetstream.subject.email"),
			AckPolicy:     jetstream.AckExplicitPolicy,
		})

		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create or update email consumer")
		}

		_, err = emailCons.Consume(func(msg jetstream.Msg) {
			logger.Info().
				Str("subject", msg.Subject()).
				Msgf("Received message: %s", string(msg.Data()))
			go func() {
				sh.EmailHandler(msg)
				msg.Ack()
			}()
		})

		if err != nil {
			logger.Error().Err(err).Msg("Failed to consume email consumer")
			return
		}

		if s.ConnectionReady != nil {
			s.ConnectionReady <- true
		}

		logger.Info().Msg("NATS subscriber is running")
		<-ctx.Done()
		logger.Info().Msg("Shutting down NATS subscriber")
	})
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
}
