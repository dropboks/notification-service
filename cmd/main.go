package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dropboks/notification-service/cmd/bootstrap"
	"github.com/dropboks/notification-service/cmd/server"
)

func main() {
	container := bootstrap.Run()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	subscriberReady := make(chan bool)
	subscriberDone := make(chan struct{})
	subscriber := &server.Subscriber{
		Container:       container,
		ConnectionReady: subscriberReady,
	}
	go func() {
		subscriber.Run(ctx)
		close(subscriberDone)
	}()
	<-subscriberReady

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	<-subscriberDone
}
