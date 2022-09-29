package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/kafka-reader/internal/app"
	"github.com/kafka-reader/internal/config"
	"github.com/kafka-reader/internal/service"
	"github.com/segmentio/kafka-go"
	"log"
)

func main() {
	log.Println("Dependencies initialization ...")

	appConfig := readAppConfig()
	kafkaReaderConfig := readKafkaReaderConfig()
	kafkaReader := initKafkaReader(kafkaReaderConfig)
	defer func(kafkaReader *kafka.Reader) {
		err := kafkaReader.Close()
		if err != nil {
			log.Fatalf("failed to close kafka reader: %v", err)
		}
	}(kafkaReader)

	eventService := service.NewEventService(kafkaReader)

	application := app.NewApp(appConfig, eventService)
	application.StartApp()
}

func readAppConfig() *config.AppConfig {
	cfg := &config.AppConfig{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to read app configs: %v", err)
	}
	log.Printf("app config read: %#v", cfg)

	return cfg
}

func readKafkaReaderConfig() *config.KafkaReaderConfig {
	cfg := &config.KafkaReaderConfig{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("failed to read kafka reader configs: %v", err)
	}
	log.Printf("kafka config read: %#v", cfg)

	return cfg
}

func initKafkaReader(config *config.KafkaReaderConfig) *kafka.Reader {
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		GroupID:        config.GroupId,
		Topic:          config.Topic,
		MinBytes:       config.MinBytes,
		MaxBytes:       config.MaxBytes,
		CommitInterval: config.CommitInterval,
	})

	return kafkaReader
}
