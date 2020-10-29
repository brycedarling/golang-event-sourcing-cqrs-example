package eventstore

// Reader ...
type Reader interface {
	Read(streamName string, position int, batchSize int) (Events, error)
	ReadAll(streamName string) (Events, error)
	ReadLast(streamName string) (*Event, error)
}
