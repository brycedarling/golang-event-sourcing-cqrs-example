package messagedb

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
)

// NewReader ...
func NewReader(db *sql.DB) eventstore.Reader {
	return &reader{db}
}

type reader struct {
	db *sql.DB
}

var _ eventstore.Reader = (*reader)(nil)

const (
	categoryMessagesSQL string = "SELECT * FROM get_category_messages($1, $2, $3)"
	streamMessagesSQL   string = "SELECT * FROM get_stream_messages($1, $2, $3)"
)

// Read ...
func (r *reader) Read(streamName string, position int, blockSize int) (events eventstore.Events, err error) {
	var query string
	if strings.Contains(streamName, "-") {
		// Entity streams have a dash
		query = streamMessagesSQL
	} else {
		// Category streams do not have a dash
		query = categoryMessagesSQL
	}

	rows, err := r.db.Query(query, streamName, position, blockSize)
	if err != nil {
		return events, err
	}
	defer rows.Close()

	for rows.Next() {
		event, err := deserializeEvent(rows)
		if err != nil {
			return events, err
		}
		events = append(events, event)
	}
	return events, nil
}

const lastStreamMessageSQL string = "SELECT * FROM get_last_stream_message($1)"

// ReadLastMessage ...
func (r *reader) ReadLast(streamName string) (*eventstore.Event, error) {
	return deserializeEvent(r.db.QueryRow(lastStreamMessageSQL, streamName))
}

// Scanner ...
type Scanner interface {
	Scan(...interface{}) error
}

func deserializeEvent(row Scanner) (*eventstore.Event, error) {
	msg := &eventstore.Event{}
	var (
		data     []byte
		metadata []byte
	)
	err := row.Scan(&msg.ID, &msg.StreamName, &msg.Type, &msg.Position, &msg.GlobalPosition, &data, &metadata, &msg.Time)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if len(data) > 0 {
		err = json.Unmarshal(data, &msg.Data)
		if err != nil {
			return nil, err
		}
	}
	if len(metadata) > 0 {
		err = json.Unmarshal(metadata, &msg.Metadata)
		if err != nil {
			return nil, err
		}
	}
	return msg, nil
}

const blockSize int = 1000

func (r *reader) ReadAll(streamName string) (events eventstore.Events, err error) {
	position := 0
	var more eventstore.Events
	for {
		more, err = r.Read(streamName, position, blockSize)
		if err != nil {
			return events, err
		}
		events = append(events, more...)
		if len(more) != blockSize {
			return events, nil
		}
		position += blockSize
	}
}
