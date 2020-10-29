package messagedb

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
)

// VersionConflictError ...
type VersionConflictError struct {
	streamName      string
	actualVersion   int
	expectedVersion *int
}

var versionConflictErrorRegex = regexp.MustCompile(".*Wrong.*Stream Version: (?P<ActualVersion>\\d+)\\)$")

// NewVersionConflictError ...
func NewVersionConflictError(event *eventstore.Event, err error) *VersionConflictError {
	errorMatches := versionConflictErrorRegex.FindStringSubmatch(err.Error())
	if len(errorMatches) == 0 {
		return nil
	}
	actualVersion, err := strconv.Atoi(errorMatches[1])
	if err != nil {
		actualVersion = -1
	}
	return &VersionConflictError{event.StreamName, actualVersion, event.ExpectedVersion}
}

func (e *VersionConflictError) Error() string {
	var expectedVersion string
	if e.expectedVersion != nil {
		expectedVersion = fmt.Sprintf("%d", *e.expectedVersion)
	} else {
		expectedVersion = fmt.Sprintf("%v", nil)
	}
	return fmt.Sprintf("Version conflict on stream %s. Expected version %s, actual %d", e.streamName, expectedVersion, e.actualVersion)
}
