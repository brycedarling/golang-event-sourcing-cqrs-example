package command

import (
	"errors"
	"strings"

	"github.com/brycedarling/go-practical-microservices/internal/domain/identity"
	"github.com/brycedarling/go-practical-microservices/internal/domain/identity/event"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
)

// RegisterCommand ...
type RegisterCommand interface {
	Execute() error
}

// NewRegisterCommand ...
func NewRegisterCommand(s eventstore.Store, iq identity.Query, ph identity.PasswordHasher,
	traceID, email, password string) (RegisterCommand, error) {
	cmd := &registerCommand{s, iq, ph, traceID, email, password}
	if err := cmd.validate(); err != nil {
		return nil, err
	}
	return cmd, nil
}

type registerCommand struct {
	eventStore     eventstore.Store
	identityQuery  identity.Query
	passwordHasher identity.PasswordHasher
	traceID        string
	email          string
	password       string
}

// Execute ...
func (cmd *registerCommand) Execute() error {
	if err := cmd.validate(); err != nil {
		return err
	}
	if err := cmd.ensureIdentityDoesNotExist(); err != nil {
		return err
	}
	hashedPassword, err := cmd.passwordHasher.GenerateFromPassword([]byte(cmd.password))
	if err != nil {
		return err
	}
	event, err := event.NewRegisterEvent(cmd.traceID, cmd.email, string(hashedPassword))
	if err != nil {
		return err
	}
	_, err = cmd.eventStore.Write(event)
	return err
}

func (cmd *registerCommand) validate() error {
	var valErrs []string
	if cmd.traceID == "" {
		valErrs = append(valErrs, "missing trace id")
	}
	if cmd.email == "" {
		valErrs = append(valErrs, "missing email")
	}
	if !strings.Contains(cmd.email, "@") {
		valErrs = append(valErrs, "invalid email format")
	}
	if cmd.password == "" {
		valErrs = append(valErrs, "missing password")
	} else if len(cmd.password) < 8 {
		valErrs = append(valErrs, "password must be at least 8 characters")
	}
	if len(valErrs) > 0 {
		return errors.New(strings.Join(valErrs, ", "))
	}
	return nil
}

// uses eventually consistent view data so it could potentially be out of date and wrong
func (cmd *registerCommand) ensureIdentityDoesNotExist() error {
	_, err := cmd.identityQuery.FindByEmail(cmd.email)
	if err == identity.ErrIdentityNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	return identity.ErrIdentityAlreadyExists
}
