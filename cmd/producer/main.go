package main

import (
	"bufio"
	"os"

	"github.com/ONSdigital/dp-notification-sender/config"
	"github.com/ONSdigital/dp-notification-sender/event"
	"github.com/ONSdigital/dp-notification-sender/schema"
	"github.com/ONSdigital/go-ns/kafka"
	"github.com/ONSdigital/go-ns/log"
)

func main() {
	log.Namespace = "notification-sender"

	config, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	// Avoid logging the neo4j FileURL as it may contain a password
	log.Debug("loaded config", log.Data{"config": config})

	kafkaBrokers := config.KafkaAddr

	kafkaProducer, err := kafka.NewProducer(kafkaBrokers, config.FilterConsumerTopic, 0)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		filterID := scanner.Text()

		log.Debug("Sending filter output event", log.Data{"filter_ouput_id": filterID})

		event := event.FilterCompleted{
			FilterID: filterID,
		}

		bytes, err := schema.FilterCompletedEvent.Marshal(event)
		if err != nil {
			log.Error(err, nil)
			os.Exit(1)
		}

		kafkaProducer.Output() <- bytes
	}
}
