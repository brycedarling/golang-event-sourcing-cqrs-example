package command

import (
	"errors"
	"strings"

	"github.com/brycedarling/go-practical-microservices/internal/domain/identity"
	"github.com/brycedarling/go-practical-microservices/internal/domain/identity/event"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
)

// AuthenticateCommand ...
type AuthenticateCommand interface {
	Execute() (*identity.Identity, error)
}

// NewAuthenticateCommand ...
func NewAuthenticateCommand(s eventstore.Store, iq identity.Query, ph identity.PasswordHasher,
	traceID, email, password string) (AuthenticateCommand, error) {
	cmd := &authenticateCommand{s, iq, ph, traceID, email, password}
	if err := cmd.validate(); err != nil {
		return nil, err
	}
	return cmd, nil
}

type authenticateCommand struct {
	eventStore     eventstore.Store
	identityQuery  identity.Query
	passwordHasher identity.PasswordHasher
	traceID        string
	email          string
	password       string
}

// Execute ...
func (cmd *authenticateCommand) Execute() (*identity.Identity, error) {
	if err := cmd.validate(); err != nil {
		return nil, err
	}
	id, err := cmd.ensureIdentityFound(cmd.loadIdentity())
	if err != nil {
		return nil, err
	}
	if err = cmd.validatePassword(id); err != nil {
		if err := cmd.writeUnauthenticatedEvent(cmd.traceID, id.UserID, err); err != nil {
			return nil, err
		}
		return nil, err
	}
	if err = cmd.writeAuthenticatedEvent(cmd.traceID, id.UserID); err != nil {
		return nil, err
	}
	return id, nil
}

func (cmd *authenticateCommand) validate() error {
	var valErrs []string
	if cmd.traceID == "" {
		valErrs = append(valErrs, "missing trace id")
	}
	if cmd.email == "" {
		valErrs = append(valErrs, "missing email")
	}
	if cmd.password == "" {
		valErrs = append(valErrs, "missing password")
	}
	if len(valErrs) > 0 {
		return errors.New(strings.Join(valErrs, ", "))
	}
	return nil
}

func (cmd *authenticateCommand) loadIdentity() (*identity.Identity, error) {
	return cmd.identityQuery.FindByEmail(cmd.email)
}

func (*authenticateCommand) ensureIdentityFound(id *identity.Identity, err error) (*identity.Identity, error) {
	if err == identity.ErrIdentityNotFound || id == nil {
		return nil, ErrAuthenticationFailed{err}
	}
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (cmd *authenticateCommand) validatePassword(id *identity.Identity) error {
	err := cmd.passwordHasher.CompareHashAndPassword([]byte(id.HashedPassword), []byte(cmd.password))
	if err != nil {
		return ErrAuthenticationFailed{err}
	}
	return nil
}

func (cmd *authenticateCommand) writeUnauthenticatedEvent(traceID, userID string, authErr error) error {
	event, err := event.NewUnauthenticatedEvent(traceID, userID, authErr)
	if err != nil {
		return err
	}
	_, err = cmd.eventStore.Write(event)
	if err != nil {
		return err
	}
	return nil
}

func (cmd *authenticateCommand) writeAuthenticatedEvent(traceID, userID string) error {
	event, err := event.NewAuthenticatedEvent(traceID, userID)
	if err != nil {
		return err
	}
	_, err = cmd.eventStore.Write(event)
	if err != nil {
		return err
	}
	return nil
}

// ErrAuthenticationFailed ...
type ErrAuthenticationFailed struct {
	error
}

func (e ErrAuthenticationFailed) Error() string {
	return e.error.Error()
}
