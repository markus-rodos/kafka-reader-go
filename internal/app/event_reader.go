package app

import (
	"github.com/kafka-reader/internal/domain"
	"github.com/kafka-reader/internal/service"
	"log"
	"sync/atomic"
	"time"
)

type EventReader struct {
	keepRunning  bool
	eventChannel chan *domain.EventHolder
	eventService *service.EventService
}

func NewEventReader(
	eventService *service.EventService,
	eventChannel chan *domain.EventHolder) *EventReader {
	return &EventReader{
		eventService: eventService,
		eventChannel: eventChannel,
		keepRunning:  true}
}

func (p *EventReader) ReadEvents() {
	log.Println("Start reading events")
	var count int64
	var start = time.Now()

	defer close(p.eventChannel)
	for p.keepRunning {
		p.eventChannel <- p.eventService.GetEvent()
		cur := atomic.AddInt64(&count, 1)
		if cur%1000 == 0 {
			log.Printf("Produced %d events at speed %.2f/s", cur, float64(cur)/time.Since(start).Seconds())
		}
	}
	log.Println("Exiting events reader")
}

func (p *EventReader) Shutdown() {
	p.keepRunning = false
}
