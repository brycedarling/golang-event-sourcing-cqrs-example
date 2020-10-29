package viewing

import "errors"

// Query ...
type Query interface {
	Initialize() error
	Find() (*Viewing, error)
	IncrementVideosWatched(globalPosition int) error
}

// ErrViewingNotFound ...
var ErrViewingNotFound = errors.New("viewing not found")
