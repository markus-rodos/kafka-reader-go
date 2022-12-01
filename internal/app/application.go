package app

import (
	"context"
	"github.com/kafka-reader/internal/config"
	"github.com/kafka-reader/internal/domain"
	"github.com/kafka-reader/internal/service"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	config       *config.AppConfig
	eventService *service.EventService
}

func NewApp(
	config *config.AppConfig,
	eventService *service.EventService) *App {

	return &App{config: config, eventService: eventService}
}

func (app *App) StartApp() {
	log.Println("Application starting ...")

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(ctx)
	defer stop()
	defer cancel()

	wg := sync.WaitGroup{}
	eventChannel := make(chan *domain.EventHolder, app.config.EventChannelSize)
	batchChannel := make(chan []*domain.Event, app.config.BatchChannelSize)
	eventsConsumer := NewEventConsumer(app.config, eventChannel, batchChannel, app.eventService)

	// handler
	for i := 0; i < app.config.WorkerCount; i++ {
		wg.Add(1)
		go eventsConsumer.HandleEvents(ctx, &wg, cancel)
	}

	// batcher
	wg.Add(1)
	go eventsConsumer.BatchEvents(ctx, &wg, cancel)

	// reader
	eventsReader := NewEventReader(app.eventService, eventChannel)
	go eventsReader.ReadEvents()

	wg.Wait()
	eventsReader.Shutdown()

	log.Println("Application stopping ...")
}
