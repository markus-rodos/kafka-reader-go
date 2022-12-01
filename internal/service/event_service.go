package service

import (
	"context"
	"encoding/json"
	"github.com/kafka-reader/internal/domain"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type EventService struct {
	Reader *kafka.Reader
}

func NewEventService(reader *kafka.Reader) *EventService {
	return &EventService{Reader: reader}
}

func (s *EventService) GetEvent() *domain.EventHolder {
	msg, err := s.Reader.ReadMessage(context.TODO())
	if err != nil {
		return &domain.EventHolder{
			Err: &domain.AppError{
				Fatal:   true,
				Message: "failed to read msg",
				Cause:   err,
			}}
	}
	val := string(msg.Value)
	event := domain.Event{}
	err = json.Unmarshal([]byte(val), &event)
	if err != nil {
		return &domain.EventHolder{
			Err: &domain.AppError{
				Message: "failed to parse msg",
				Payload: val,
				Cause:   err,
			}}
	}

	return &domain.EventHolder{Event: &event}
}

func (s *EventService) Process(event *domain.Event) error {
	// do some event processing
	log.Printf("Processing event: %+v", event)
	time.Sleep(200 * time.Microsecond)

	return nil
}
