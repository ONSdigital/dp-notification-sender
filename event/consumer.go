package event

import (
	"context"
	errs "errors"
	"net/smtp"

	"github.com/ONSdigital/dp-notification-sender/schema"
	"github.com/ONSdigital/go-ns/kafka"
	"github.com/ONSdigital/go-ns/log"
)

//go:generate moq -out eventtest/handler.go -pkg eventtest . Handler

// MessageConsumer provides a generic interface for consuming []byte messages
type MessageConsumer interface {
	Incoming() chan kafka.Message
}

// Handler represents a handler for processing a single event.
type Handler interface {
	Handle(filterCompletedEvent *FilterCompleted) error
}

// Consumer consumes event messages.
type Consumer struct {
	closing chan bool
	closed  chan bool
}

// NewConsumer returns a new consumer instance.
func NewConsumer() *Consumer {
	return &Consumer{
		closing: make(chan bool),
		closed:  make(chan bool),
	}
}

// Consume converts messages to event instances, and pass the event to the provided handler.
func (consumer *Consumer) Consume(messageConsumer MessageConsumer, isReady chan bool) {

	go func() {
		defer close(consumer.closed)

		log.Info("waiting for authentication before consuming messages", nil)
		if <-isReady {
			log.Info("starting to consume messages", nil)
			for {
				select {
				case message := <-messageConsumer.Incoming():

					processMessage(message)

				case <-consumer.closing:
					log.Info("closing event consumer loop", nil)
					return
				}
			}
		}
	}()

}

// Close safely closes the consumer and releases all resources
func (consumer *Consumer) Close(ctx context.Context) (err error) {

	if ctx == nil {
		ctx = context.Background()
	}

	close(consumer.closing)

	select {
	case <-consumer.closed:
		log.Info("successfully closed event consumer", nil)
		return nil
	case <-ctx.Done():
		log.Info("shutdown context time exceeded, skipping graceful shutdown of event consumer", nil)
		return errs.New("Shutdown context timed out")
	}
}

func processMessage(message kafka.Message) {
	event, err := unmarshal(message)
	if err != nil {
		log.Error(err, log.Data{"message": "failed to unmarshal event"})
		return
	}

	log.Debug("event received", log.Data{"event": event})

	err = smtp.SendMail(
		"localhost:1025",
		nil,
		"no-reply@ons.gov.uk",
		[]string{event.Email},
		[]byte(`To: Bob <bob@email.com>
From: ONS <ons@ons.gov.uk>
Subject: Some kind of subject

This is the email body.`),
	)
	if err != nil {
		log.Error(err, log.Data{"message": "failed to unmarshal event"})
		return
	}

	log.Debug("event processed - committing message", log.Data{"event": event})
	message.Commit()
	log.Debug("message committed", log.Data{"event": event})

}

// unmarshal converts a event instance to []byte.
func unmarshal(message kafka.Message) (*FilterCompleted, error) {
	var event FilterCompleted
	err := schema.FilterCompletedEvent.Unmarshal(message.GetData(), &event)
	return &event, err
}
