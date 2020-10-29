package identity

// PasswordHasher ...
type PasswordHasher interface {
	GenerateFromPassword(password []byte) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}
