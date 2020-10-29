package identity

import "errors"

// Query ...
type Query interface {
	CreateIdentity(*Identity) error
	FindByEmail(email string) (*Identity, error)
}

// ErrIdentityAlreadyExists ...
var ErrIdentityAlreadyExists = errors.New("identity already exists")

// ErrIdentityNotFound ...
var ErrIdentityNotFound = errors.New("identity not found")
