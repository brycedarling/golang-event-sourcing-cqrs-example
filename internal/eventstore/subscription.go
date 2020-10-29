package eventstore

// Subscription ...
type Subscription interface {
	Subscribe()
	Unsubscribe()
}

// Subscriber ...
type Subscriber func(*Event)

// Subscribers ...
type Subscribers map[string]Subscriber
