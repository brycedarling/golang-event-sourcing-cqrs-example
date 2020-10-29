package identity

import (
	"github.com/brycedarling/go-practical-microservices/internal/application"
	"github.com/brycedarling/go-practical-microservices/internal/domain/identity/event"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
)

// Aggregator ...
type Aggregator struct {
	subscription eventstore.Subscription
}

var _ application.Aggregator = (*Aggregator)(nil)

// NewAggregator ...
func NewAggregator(conf *config.Config, h RegisteredEventHandler) *Aggregator {
	return &Aggregator{conf.EventStore.CreateSubscription(
		event.StreamIdentity, "aggregator:identity", eventstore.Subscribers{
			event.TypeRegistered: h.Registered,
		})}
}

// Start ...
func (a *Aggregator) Start() {
	a.subscription.Subscribe()
}

// Stop ...
func (a *Aggregator) Stop() {
	a.subscription.Unsubscribe()
}
