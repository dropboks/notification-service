package service

import (
	"github.com/dropboks/notification-service/internal/infrastructure/mail"
	mq "github.com/dropboks/notification-service/internal/infrastructure/message-queue"
	"github.com/dropboks/sharedlib/dto"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type (
	SubscriberService interface {
		SendEmail(msg dto.MailNotificationMessage) error
	}
	subscriberService struct {
		logger       zerolog.Logger
		natsInstance mq.Nats
		mail         mail.Mail
	}
)

func NewSubscriberService(natsIntance mq.Nats, logger zerolog.Logger, mail mail.Mail) SubscriberService {
	return &subscriberService{
		logger:       logger,
		natsInstance: natsIntance,
		mail:         mail,
	}
}

func (s *subscriberService) SendEmail(msg dto.MailNotificationMessage) error {
	if msg.MsgType == "welcome" {
		s.mail.SetSender(viper.GetString("mail.sender"))
		s.mail.SetSubject("Welcome to Dropboks!!")
		s.mail.SetReciever(msg.Receiver...)
		if err := s.mail.SetBody("welcome.html", struct {
			Email string
		}{
			Email: msg.Receiver[0],
		}); err != nil {
			s.logger.Error().Err(err).Msg("error set body html")
			return err
		}
		if err := s.mail.Send(); err != nil {
			s.logger.Error().Err(err).Msg("error send email")
			return err
		}
	}
	return nil
}
