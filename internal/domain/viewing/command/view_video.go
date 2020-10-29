package command

import (
	"errors"
	"strings"

	"github.com/brycedarling/go-practical-microservices/internal/domain/viewing/event"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
)

// ViewVideoCommand ...
type ViewVideoCommand interface {
	Execute() error
}

// NewViewVideoCommand ...
func NewViewVideoCommand(s eventstore.Store, traceID string, userID *string, videoID string) (ViewVideoCommand, error) {
	cmd := &viewVideoCommand{s, traceID, userID, videoID}
	if err := cmd.validate(); err != nil {
		return nil, err
	}
	return cmd, nil
}

type viewVideoCommand struct {
	eventStore eventstore.Store
	TraceID    string
	UserID     *string
	VideoID    string
}

func (cmd *viewVideoCommand) validate() error {
	var valErrs []string
	if cmd.TraceID == "" {
		valErrs = append(valErrs, "missing trace id")
	}
	if cmd.UserID == nil || *cmd.UserID == "" {
		valErrs = append(valErrs, "missing user id")
	}
	if cmd.VideoID == "" {
		valErrs = append(valErrs, "missing video id")
	}
	if len(valErrs) > 0 {
		return errors.New(strings.Join(valErrs, ", "))
	}
	return nil
}

func (cmd *viewVideoCommand) Execute() error {
	event, err := event.NewVideoViewedEvent(cmd.TraceID, *cmd.UserID, cmd.VideoID)
	if err != nil {
		return err
	}
	_, err = cmd.eventStore.Write(event)
	return err
}
