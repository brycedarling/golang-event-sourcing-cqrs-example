package eventstore

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Events ...
type Events []*Event

// Event ...
type Event struct {
	ID              string
	StreamName      string
	Type            string
	Data            map[string]interface{}
	Metadata        map[string]interface{}
	ExpectedVersion *int
	Position        int
	GlobalPosition  int
	Time            time.Time
}

// NewEvent ...
func NewEvent(streamName, eventType string) (*Event, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	event := &Event{
		ID:         id.String(),
		StreamName: streamName,
		Type:       eventType,
	}
	if err := event.validate(); err != nil {
		return nil, err
	}
	return event, nil
}

// ErrIDRequired ...
var ErrIDRequired = errors.New("Event must have an id")

// ErrStreamNameRequired ...
var ErrStreamNameRequired = errors.New("Event must have a stream name")

// ErrTypeRequired ...
var ErrTypeRequired = errors.New("Event must have a type")

func (event *Event) validate() error {
	if len(event.ID) == 0 {
		return ErrIDRequired
	}
	if len(event.StreamName) == 0 {
		return ErrStreamNameRequired
	}
	if len(event.Type) == 0 {
		return ErrTypeRequired
	}
	return nil
}
