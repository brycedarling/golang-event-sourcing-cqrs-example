package event

import (
	"errors"
	"fmt"

	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
)

// StreamAuthentication ...
const StreamAuthentication string = "authentication"

// TypeAuthenticated ...
const TypeAuthenticated string = "Authenticated"

// NewAuthenticatedEvent ...
func NewAuthenticatedEvent(traceID, userID string) (*eventstore.Event, error) {
	if err := (&authenticatedEvent{traceID, userID}).validate(); err != nil {
		return nil, err
	}
	streamName := fmt.Sprintf("%s-%s", StreamAuthentication, userID)
	event, err := eventstore.NewEvent(streamName, TypeAuthenticated)
	if err != nil {
		return nil, err
	}
	event.Data = map[string]interface{}{
		"userID": userID,
	}
	event.Metadata = map[string]interface{}{
		"traceID": traceID,
		"userID":  userID,
	}
	return event, nil
}

type authenticatedEvent struct {
	traceID string
	userID  string
}

func (event *authenticatedEvent) validate() error {
	if event.traceID == "" {
		return errors.New("missing trace id")
	}
	if event.userID == "" {
		return errors.New("missing user id")
	}
	return nil
}
