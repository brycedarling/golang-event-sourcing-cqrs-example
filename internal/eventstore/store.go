package eventstore

// Store ...
type Store interface {
	Write(*Event) (int, error)
	CreateSubscription(streamName, subscriberID string, subscribers Subscribers) Subscription
	Read(streamName string, position int, batchSize int) (Events, error)
	ReadAll(streamName string) (Events, error)
	ReadLast(streamName string) (*Event, error)
}
