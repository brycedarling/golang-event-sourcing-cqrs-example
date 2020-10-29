package identity

import (
	"log"

	"github.com/brycedarling/go-practical-microservices/internal/domain/identity"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
)

// RegisteredEventHandler ...
type RegisteredEventHandler interface {
	Registered(*eventstore.Event)
}

// NewRegisteredEventHandler ...
func NewRegisteredEventHandler(conf *config.Config) RegisteredEventHandler {
	return &registeredEventHandler{conf.IdentityQuery}
}

type registeredEventHandler struct {
	identityQuery identity.Query
}

func (h *registeredEventHandler) Registered(event *eventstore.Event) {
	userID := event.Data["userID"].(string)
	email := event.Data["email"].(string)
	hashedPassword := event.Data["hashedPassword"].(string)
	id := identity.NewIdentity(userID, email, hashedPassword)
	err := h.identityQuery.CreateIdentity(id)
	if err != nil && err != identity.ErrIdentityAlreadyExists {
		log.Println("Unexpected error creating identity:", err)
	}
}
