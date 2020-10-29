package messagedb

import (
	"fmt"
	"log"
	"time"

	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore/event"
)

// NewSubscription ...
func NewSubscription(store eventstore.Store, streamName, subscriberID string, subscribers eventstore.Subscribers) eventstore.Subscription {
	return &subscription{
		eventStore:                     store,
		streamName:                     streamName,
		subscriberID:                   subscriberID,
		subscriberStreamName:           fmt.Sprintf("subscriberPosition-%s", subscriberID),
		subscribers:                    subscribers,
		currentPosition:                0,
		messagesSinceLastPositionWrite: 0,
		isPolling:                      false,
		positionUpdateInterval:         99,
		messagesPerTick:                100,
		tickIntervalMS:                 100 * time.Millisecond,
	}
}

type subscription struct {
	eventStore                     eventstore.Store
	streamName                     string
	subscriberID                   string
	subscriberStreamName           string
	subscribers                    map[string]eventstore.Subscriber
	currentPosition                int
	messagesSinceLastPositionWrite int
	isPolling                      bool
	positionUpdateInterval         int
	messagesPerTick                int
	tickIntervalMS                 time.Duration
}

var _ eventstore.Subscription = (*subscription)(nil)

// Start ...
func (s *subscription) Subscribe() {
	log.Printf("Subscribing to %s as %s", s.streamName, s.subscriberID)

	s.loadPosition()

	s.poll()
}

// Stop ...
func (s *subscription) Unsubscribe() {
	log.Printf("Unsubscribing from %s as %s", s.streamName, s.subscriberID)

	s.isPolling = false
}

func (s *subscription) loadPosition() {
	msg, err := s.eventStore.ReadLast(s.subscriberStreamName)
	if err != nil {
		log.Println("Error loading position", err)
		return
	}
	if msg == nil {
		return
	}
	if position, ok := msg.Data["position"].(float64); ok {
		s.currentPosition = int(position)
	}
}

func (s *subscription) poll() {
	s.isPolling = true

	ticker := time.NewTicker(s.tickIntervalMS)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if !s.isPolling {
					close(quit)
				}
				s.tick()
			case <-quit:
				ticker.Stop()
			}
		}
	}()
}

func (s *subscription) tick() {
	msgs, err := s.nextBatchOfMessages()
	if err != nil {
		log.Printf("Error fetching batch: %s", err)
		s.Unsubscribe()
	}
	err = s.processBatch(msgs)
	if err != nil {
		log.Printf("Error processing batch: %s", err)
		s.Unsubscribe()
	}
}

func (s *subscription) nextBatchOfMessages() (eventstore.Events, error) {
	return s.eventStore.Read(s.streamName, s.currentPosition+1, s.messagesPerTick)
}

func (s *subscription) processBatch(events eventstore.Events) error {
	for _, event := range events {
		subscriber, ok := s.subscribers[event.Type]
		if !ok {
			continue
		}

		subscriber(event)

		err := s.updateReadPosition(event.GlobalPosition)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *subscription) updateReadPosition(position int) error {
	s.currentPosition = position
	s.messagesSinceLastPositionWrite++

	if s.messagesSinceLastPositionWrite == s.positionUpdateInterval {
		s.messagesSinceLastPositionWrite = 0

		event, err := event.NewReadEvent(s.subscriberStreamName, position)
		if err != nil {
			return err
		}
		_, err = s.eventStore.Write(event)
		return err
	}

	return nil
}
