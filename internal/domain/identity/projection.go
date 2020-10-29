package identity

import (
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
)

// Projection ...
type Projection *Identity

// Projector ...
type Projector interface {
	Project(eventstore.Events) Projection
}

// NewProjector ...
func NewProjector() Projector {
	return projector{
		"Registered": projectRegisteredEvent,
	}
}

type projector map[string]func(Projection, *eventstore.Event) Projection

func (p projector) Project(events eventstore.Events) Projection {
	identity := &Identity{}
	for _, event := range events {
		if project, ok := p[event.Type]; ok {
			identity = project(identity, event)
		}
	}
	return identity
}

func projectRegisteredEvent(identity Projection, event *eventstore.Event) Projection {
	identity.UserID = event.Data["userID"].(string)
	identity.Email = event.Data["email"].(string)
	identity.HashedPassword = event.Data["hashedPassword"].(string)
	identity.IsRegistered = true
	return identity
}
