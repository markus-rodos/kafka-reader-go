package app

import (
	"context"
	"github.com/kafka-reader/internal/config"
	"github.com/kafka-reader/internal/domain"
	"github.com/kafka-reader/internal/service"
	"github.com/kafka-reader/internal/util"
	"log"
	"sync"
	"time"
)

type EventConsumer struct {
	config       *config.AppConfig
	buffer       []*domain.Event
	eventChannel chan *domain.EventHolder
	batchChannel chan []*domain.Event
	eventService *service.EventService
}

func NewEventConsumer(
	config *config.AppConfig,
	eventChannel chan *domain.EventHolder,
	batchChannel chan []*domain.Event,
	eventService *service.EventService) *EventConsumer {
	return &EventConsumer{
		config:       config,
		buffer:       make([]*domain.Event, 0, config.BufferSize),
		eventChannel: eventChannel,
		batchChannel: batchChannel,
		eventService: eventService}
}

func (c *EventConsumer) HandleEvents(
	ctx context.Context,
	wg *sync.WaitGroup,
	cancel func()) {
	log.Println("Start handling events")
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("An error occurred. Exiting events handler")
			return
		case events := <-c.batchChannel:
			for _, event := range events {
				err := c.eventService.Process(event)
				if err != nil {
					log.Printf("failed to process event: %v", err)
				}
				if err != nil && util.IsFatal(err) {
					cancel()
					return
				}
			}
		}
	}
}

func (c *EventConsumer) BatchEvents(ctx context.Context, wg *sync.WaitGroup, cancel func()) {
	ticker := time.NewTicker(time.Duration(c.config.TickerIntervalSeconds) * time.Second)
	defer ticker.Stop()

	log.Println("Start batching events")
	defer wg.Done()
	defer close(c.batchChannel)
	for {
		select {
		case <-ctx.Done():
			log.Println("An error occurred. Exiting events batcher")
			return
		case <-ticker.C: // in case the producer is slow
			log.Printf("Ticker trigerred at %v", time.Now().UTC().Format(time.RFC3339))
			c.flushBuffer()
		case eventHolder := <-c.eventChannel:
			if eventHolder.Event != nil {
				c.appendAndFlush(eventHolder.Event)
			}
			if eventHolder.Err != nil {
				log.Printf("Failed to get message from source: %v", eventHolder.Err)
			}
			if eventHolder.Err != nil && util.IsFatal(eventHolder.Err) {
				cancel()
				return
			}
		}
	}
}

func (c *EventConsumer) flushBuffer() {
	if len(c.buffer) > 0 {
		c.batchChannel <- c.buffer
		c.buffer = make([]*domain.Event, 0, c.config.BufferSize)
	}
}

func (c *EventConsumer) appendAndFlush(event *domain.Event) {
	c.buffer = append(c.buffer, event)
	if len(c.buffer) >= c.config.BufferSize {
		c.flushBuffer()
	}
}
