package messagedb

import (
	"database/sql"

	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
)

// NewMessageDB ...
func NewMessageDB(db *sql.DB) eventstore.Store {
	return &messageDB{
		writer: NewWriter(db),
		reader: NewReader(db),
	}
}

type messageDB struct {
	writer eventstore.Writer
	reader eventstore.Reader
}

var _ eventstore.Store = (*messageDB)(nil)

// Write ...
func (db *messageDB) Write(event *eventstore.Event) (int, error) {
	return db.writer.Write(event)
}

// CreateSubscription ...
func (db *messageDB) CreateSubscription(streamName, subscriberID string, subscribers eventstore.Subscribers) eventstore.Subscription {
	return NewSubscription(db, streamName, subscriberID, subscribers)
}

// ReadLast ...
func (db *messageDB) ReadLast(streamName string) (*eventstore.Event, error) {
	event, err := db.reader.ReadLast(streamName)
	if err != nil {
		return event, err
	}
	return event, err
}

// Read ...
func (db *messageDB) Read(streamName string, position, batchSize int) (eventstore.Events, error) {
	return db.reader.Read(streamName, position, batchSize)
}

// ReadAll ...
func (db *messageDB) ReadAll(streamName string) (eventstore.Events, error) {
	return db.reader.ReadAll(streamName)
}
