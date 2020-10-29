package identity

import (
	"fmt"
)

// Identity ...
type Identity struct {
	UserID         string `redis:"user_id" json:"userId"`
	Email          string `redis:"email" json:"email"`
	HashedPassword string `redis:"hashed_password" json:"-"`
	IsRegistered   bool   `redis:"-" json:"-"`
}

// NewIdentity ...
func NewIdentity(userID, email, hashedPassword string) *Identity {
	return &Identity{userID, email, hashedPassword, false}
}

const identityFormat string = "<Identity UserID=\"%v\" Email=\"%v\" HashedPassword=\"%v\" IsRegistered=%v>"

func (i *Identity) String() string {
	return fmt.Sprintf(identityFormat, i.UserID, i.Email, i.HashedPassword, i.IsRegistered)
}
