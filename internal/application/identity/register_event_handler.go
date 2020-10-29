package identity

import (
	"fmt"
	"log"

	"github.com/brycedarling/go-practical-microservices/internal/domain/identity"
	"github.com/brycedarling/go-practical-microservices/internal/domain/identity/event"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
)

// RegisterEventHandler ...
type RegisterEventHandler interface {
	Register(*eventstore.Event)
}

// NewRegisterEventHandler ...
func NewRegisterEventHandler(conf *config.Config, ip identity.Projector) RegisterEventHandler {
	return &registerEventHandler{conf.EventStore, ip}
}

type registerEventHandler struct {
	eventStore        eventstore.Store
	identityProjector identity.Projector
}

func (h *registerEventHandler) Register(cmd *eventstore.Event) {
	identity, err := h.loadIdentity(cmd.Data["userID"].(string))
	if err != nil {
		log.Printf("Error loading identity: %s", err)
	}
	if identity.IsRegistered {
		return
	}

	traceID := cmd.Metadata["traceID"].(string)
	userID := cmd.Data["userID"].(string)
	email := cmd.Data["email"].(string)
	hashedPassword := cmd.Data["hashedPassword"].(string)

	event, err := event.NewRegisteredEvent(traceID, userID, email, hashedPassword)
	if err != nil {
		log.Printf("Error creating registered event: %s", err)
		return
	}
	_, err = h.eventStore.Write(event)
	if err != nil {
		log.Printf("Error writing registered event: %s", err)
	}
}

func (h *registerEventHandler) loadIdentity(userID string) (*identity.Identity, error) {
	identityStreamName := fmt.Sprintf("identity-%s", userID)
	msgs, err := h.eventStore.ReadAll(identityStreamName)
	if err != nil {
		return nil, err
	}
	return h.identityProjector.Project(msgs), nil
}
