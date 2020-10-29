package event

import (
	"errors"
	"fmt"
	"strings"

	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
)

// StreamViewing ...
const StreamViewing string = "viewing"

// TypeVideoViewed ...
const TypeVideoViewed string = "VideoViewed"

// NewVideoViewedEvent ...
func NewVideoViewedEvent(traceID, userID, videoID string) (*eventstore.Event, error) {
	if err := (&videoViewedEvent{traceID, userID, videoID}).validate(); err != nil {
		return nil, err
	}
	streamName := fmt.Sprintf("%s-%s", StreamViewing, videoID)
	e, err := eventstore.NewEvent(streamName, TypeVideoViewed)
	if err != nil {
		return nil, err
	}
	e.Data = map[string]interface{}{
		"userID":  userID,
		"videoID": videoID,
	}
	e.Metadata = map[string]interface{}{
		"traceID": traceID,
		"userID":  userID,
	}
	return e, nil
}

type videoViewedEvent struct {
	traceID string
	userID  string
	videoID string
}

func (event *videoViewedEvent) validate() error {
	var valErrs []string
	if event.traceID == "" {
		valErrs = append(valErrs, "missing trace id")
	}
	if event.userID == "" {
		valErrs = append(valErrs, "missing user id")
	}
	if event.videoID == "" {
		valErrs = append(valErrs, "missing video id")
	}
	if len(valErrs) > 0 {
		return errors.New(strings.Join(valErrs, ", "))
	}
	return nil
}
