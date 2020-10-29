package event

import (
	"errors"
	"fmt"
	"strings"

	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
)

// StreamIdentity ...
const StreamIdentity string = "identity"

// TypeRegistered ...
const TypeRegistered string = "Registered"

// NewRegisteredEvent ...
func NewRegisteredEvent(traceID, userID, email, hashedPassword string) (*eventstore.Event, error) {
	if err := (&registeredEvent{traceID, userID, email, hashedPassword}).validate(); err != nil {
		return nil, err
	}
	streamName := fmt.Sprintf("%s-%s", StreamIdentity, userID)
	event, err := eventstore.NewEvent(streamName, TypeRegistered)
	if err != nil {
		return nil, err
	}
	event.Data = map[string]interface{}{
		"userID":         userID,
		"email":          email,
		"hashedPassword": hashedPassword,
	}
	event.Metadata = map[string]interface{}{
		"traceID": traceID,
		"userID":  userID,
	}
	return event, nil
}

type registeredEvent struct {
	traceID        string
	userID         string
	email          string
	hashedPassword string
}

func (event *registeredEvent) validate() error {
	var valErrs []string
	if event.traceID == "" {
		valErrs = append(valErrs, "missing trace id")
	}
	if event.userID == "" {
		valErrs = append(valErrs, "missing user id")
	}
	if event.email == "" {
		valErrs = append(valErrs, "missing email")
	}
	if event.hashedPassword == "" {
		valErrs = append(valErrs, "missing hashed password")
	}
	if len(valErrs) > 0 {
		return errors.New(strings.Join(valErrs, ", "))
	}
	return nil
}
