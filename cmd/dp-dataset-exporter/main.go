package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-notification-sender/config"
	"github.com/ONSdigital/dp-notification-sender/event"
	"github.com/ONSdigital/go-ns/kafka"
	"github.com/ONSdigital/go-ns/log"
)

func main() {
	log.Namespace = "dp-notification-sender"
	log.Info("Starting dp-notification-sender service", nil)

	cfg, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	log.Info("loaded config", log.Data{"config": cfg})

	// a channel used to signal a graceful exit is required.
	errorChannel := make(chan error)

	kafkaBrokers := cfg.KafkaAddr
	kafkaConsumer, err := kafka.NewConsumerGroup(
		kafkaBrokers,
		cfg.FilterConsumerTopic,
		cfg.FilterConsumerGroup,
		kafka.OffsetNewest)
	exitIfError(err)

	isReady := make(chan bool)

	eventConsumer := event.NewConsumer()
	eventConsumer.Consume(kafkaConsumer, isReady)

	shutdownGracefully := func() {

		ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)

		// gracefully dispose resources
		err = eventConsumer.Close(ctx)
		logIfError(err)

		err = kafkaConsumer.Close(ctx)
		logIfError(err)

		// cancel the timer in the shutdown context.
		cancel()

		log.Info("graceful shutdown was successful", nil)
		os.Exit(0)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	isReady <- true

	for {
		select {
		case err := <-kafkaConsumer.Errors():
			log.ErrorC("kafka consumer", err, nil)
			shutdownGracefully()
		case err := <-errorChannel:
			log.ErrorC("error channel", err, nil)
			shutdownGracefully()
		case <-signals:
			log.Debug("os signal received", nil)
			shutdownGracefully()
		}
	}
}

func exitIfError(err error) {
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}
}

func logIfError(err error) {
	if err != nil {
		log.Error(err, nil)
	}
}
