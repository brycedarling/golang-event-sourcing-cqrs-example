package identity

import (
	"github.com/brycedarling/go-practical-microservices/internal/application"
	"github.com/brycedarling/go-practical-microservices/internal/domain/identity/event"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
)

// Component ...
type Component struct {
	subscription eventstore.Subscription
}

var _ application.Component = (*Component)(nil)

// NewComponent ...
func NewComponent(conf *config.Config, h RegisterEventHandler) *Component {
	return &Component{conf.EventStore.CreateSubscription(
		event.StreamIdentityCommand, "component:identity:command", eventstore.Subscribers{
			event.TypeRegister: h.Register,
		})}
}

// Start ...
func (c *Component) Start() {
	c.subscription.Subscribe()
}

// Stop ...
func (c *Component) Stop() {
	c.subscription.Unsubscribe()
}
