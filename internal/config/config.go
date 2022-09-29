package config

import "time"

type KafkaReaderConfig struct {
	Brokers        []string      `env:"KAFKA_BOOTSTRAP_SERVERS" envDefault:"localhost:9092" envSeparator:","`
	GroupId        string        `env:"KAFKA_GROUP_ID" envDefault:"kafka-reader"`
	Topic          string        `env:"KAFKA_TOPIC" envDefault:"events.all"`
	MinBytes       int           `env:"KAFKA_MIN_BYTES" envDefault:"100"`
	MaxBytes       int           `env:"KAFKA_MAX_BYTES" envDefault:"1000"`
	CommitInterval time.Duration `env:"KAFKA_COMMIT_INTERVAL" envDefault:"5s"`
}

type AppConfig struct {
	EventChannelSize      int `env:"EVENT_CHANNEL_SIZE" envDefault:"100"`
	BatchChannelSize      int `env:"BATCH_CHANNEL_SIZE" envDefault:"5"`
	WorkerCount           int `env:"WORKER_COUNT" envDefault:"5"`
	BufferSize            int `env:"BUFFER_SIZE" envDefault:"20"`
	TickerIntervalSeconds int `env:"TICKER_INTERVAL" envDefault:"60"`
}
