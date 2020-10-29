package event

import (
	"errors"
	"fmt"
	"strings"

	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/google/uuid"
)

// StreamIdentityCommand ...
const StreamIdentityCommand string = "identity:command"

// TypeRegister ...
const TypeRegister string = "Register"

// NewRegisterEvent ...
func NewRegisterEvent(traceID, email, hashedPassword string) (*eventstore.Event, error) {
	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	if err := (&registerEvent{traceID, userID.String(), email, hashedPassword}).validate(); err != nil {
		return nil, err
	}
	streamName := fmt.Sprintf("%s-%s", StreamIdentityCommand, userID)
	event, err := eventstore.NewEvent(streamName, TypeRegister)
	if err != nil {
		return nil, err
	}
	event.Data = map[string]interface{}{
		"userID":         userID.String(),
		"email":          email,
		"hashedPassword": hashedPassword,
	}
	event.Metadata = map[string]interface{}{
		"traceID": traceID,
		"userID":  userID.String(),
	}
	return event, nil
}

type registerEvent struct {
	traceID        string
	userID         string
	email          string
	hashedPassword string
}

func (event *registerEvent) validate() error {
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
