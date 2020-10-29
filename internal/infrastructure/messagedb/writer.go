package messagedb

import (
	"database/sql"
	"encoding/json"

	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/google/uuid"
)

// NewWriter ...
func NewWriter(db *sql.DB) eventstore.Writer {
	return &writer{db}
}

type writer struct {
	db *sql.DB
}

var _ eventstore.Writer = (*writer)(nil)

const writeSQL string = "SELECT message_store.write_message($1, $2, $3, $4, $5, $6)"

// Write ...
func (w *writer) Write(event *eventstore.Event) (int, error) {
	if event.ID == "" {
		id, err := uuid.NewUUID()
		if err != nil {
			return 0, err
		}
		event.ID = id.String()
	}
	data, err := json.Marshal(event.Data)
	if err != nil {
		return 0, err
	}
	metadata, err := json.Marshal(event.Metadata)
	if err != nil {
		return 0, err
	}
	res := w.db.QueryRow(writeSQL, event.ID, event.StreamName, event.Type, data, metadata, event.ExpectedVersion)
	var nextPosition int
	err = res.Scan(&nextPosition)
	if err != nil {
		versionConflictErr := NewVersionConflictError(event, err)
		if versionConflictErr != nil {
			return 0, versionConflictErr
		}
		return 0, err
	}
	return nextPosition, nil
}
