package identity

import (
	"fmt"

	"github.com/brycedarling/go-practical-microservices/internal/domain/identity"
	"github.com/gomodule/redigo/redis"
)

// NewQueryRedis ...
func NewQueryRedis(pool *redis.Pool) identity.Query {
	return &queryRedis{pool}
}

type queryRedis struct {
	pool *redis.Pool
}

var _ identity.Query = (*queryRedis)(nil)

// CreateIdentity ...
func (q *queryRedis) CreateIdentity(id *identity.Identity) error {
	conn := q.pool.Get()
	defer conn.Close()

	key := identityKey(id.Email)

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return err
	}
	if exists {
		return identity.ErrIdentityAlreadyExists
	}

	_, err = conn.Do("HMSET", key,
		"user_id", id.UserID,
		"email", id.Email,
		"hashed_password", id.HashedPassword,
	)
	return err
}

// FindByEmail ...
func (q *queryRedis) FindByEmail(email string) (*identity.Identity, error) {
	conn := q.pool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", identityKey(email)))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, identity.ErrIdentityNotFound
	}

	var id identity.Identity
	err = redis.ScanStruct(values, &id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func identityKey(email string) string {
	return fmt.Sprintf("identity:%s", email)
}
