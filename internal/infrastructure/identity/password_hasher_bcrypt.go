package identity

import (
	"github.com/brycedarling/go-practical-microservices/internal/domain/identity"
	"golang.org/x/crypto/bcrypt"
)

// NewPasswordHasherBcrypt ...
func NewPasswordHasherBcrypt() identity.PasswordHasher {
	return &passwordHasherBcrypt{}
}

type passwordHasherBcrypt struct{}

var _ identity.PasswordHasher = (*passwordHasherBcrypt)(nil)

func (*passwordHasherBcrypt) GenerateFromPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func (*passwordHasherBcrypt) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
