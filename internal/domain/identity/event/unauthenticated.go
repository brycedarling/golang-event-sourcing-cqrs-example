package event

import (
	"errors"
	"fmt"

	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
)

// TypeUnauthenticated ...
const TypeUnauthenticated string = "Unauthenticated"

// NewUnauthenticatedEvent ...
func NewUnauthenticatedEvent(traceID, userID string, authError error) (*eventstore.Event, error) {
	if err := (&unauthenticatedEvent{traceID, userID}).validate(); err != nil {
		return nil, err
	}
	streamName := fmt.Sprintf("%s-%s", StreamAuthentication, userID)
	event, err := eventstore.NewEvent(streamName, TypeUnauthenticated)
	if err != nil {
		return nil, err
	}
	event.Data = map[string]interface{}{
		"userID": userID,
		"error":  authError.Error(),
	}
	event.Metadata = map[string]interface{}{
		"traceID": traceID,
		"userID":  userID,
	}
	return event, nil
}

type unauthenticatedEvent struct {
	traceID string
	userID  string
}

func (event *unauthenticatedEvent) validate() error {
	if event.traceID == "" {
		return errors.New("missing trace id")
	}
	if event.userID == "" {
		return errors.New("missing user id")
	}
	return nil
}
