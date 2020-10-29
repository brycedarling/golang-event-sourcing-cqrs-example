package event

import (
	"errors"
	"strings"

	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
)

// TypeRead ...
const TypeRead string = "Read"

// ReadPositionKey ...
const ReadPositionKey string = "position"

// NewReadEvent ...
func NewReadEvent(streamName string, position int) (*eventstore.Event, error) {
	if err := (&readEvent{streamName, position}).validate(); err != nil {
		return nil, err
	}
	event, err := eventstore.NewEvent(streamName, TypeRead)
	if err != nil {
		return nil, err
	}
	event.Data = map[string]interface{}{
		ReadPositionKey: position,
	}
	return event, nil
}

type readEvent struct {
	streamName string
	position   int
}

func (event *readEvent) validate() error {
	var valErrs []string
	if event.streamName == "" {
		valErrs = append(valErrs, "missing stream name")
	}
	if event.position < 0 {
		valErrs = append(valErrs, "missing position")
	}
	if len(valErrs) > 0 {
		return errors.New(strings.Join(valErrs, ", "))
	}
	return nil
}
