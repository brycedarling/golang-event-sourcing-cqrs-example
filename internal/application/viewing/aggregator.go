package viewing

import (
	"log"

	"github.com/brycedarling/go-practical-microservices/internal/application"
	"github.com/brycedarling/go-practical-microservices/internal/domain/viewing"
	"github.com/brycedarling/go-practical-microservices/internal/domain/viewing/event"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
)

// Aggregator ...
type Aggregator struct {
	viewingQuery viewing.Query
	subscription eventstore.Subscription
}

var _ application.Aggregator = (*Aggregator)(nil)

// NewAggregator ...
func NewAggregator(conf *config.Config) *Aggregator {
	a := &Aggregator{conf.ViewingQuery, nil}
	a.subscription = conf.EventStore.CreateSubscription(
		event.StreamViewing, "aggregator:viewing", eventstore.Subscribers{
			event.TypeVideoViewed: a.videoViewed,
		})
	return a
}

// Start ...
func (a *Aggregator) Start() {
	a.viewingQuery.Initialize()

	a.subscription.Subscribe()
}

// Stop ...
func (a *Aggregator) Stop() {
	a.subscription.Unsubscribe()
}

func (a *Aggregator) videoViewed(event *eventstore.Event) {
	err := a.viewingQuery.IncrementVideosWatched(event.GlobalPosition)
	if err != nil {
		log.Println("Unexpected error incrementing videos watched:", err)
	}
}
