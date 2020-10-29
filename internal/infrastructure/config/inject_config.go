//+build wireinject

package config

import (
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/identity"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/messagedb"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/viewing"
	"github.com/google/wire"
)

// InitializeConfig ...
func InitializeConfig(env *Env) (*Config, func(), error) {
	wire.Build(
		NewDB,
		messagedb.NewMessageDB,
		NewRedisPool,
		identity.NewQueryRedis,
		viewing.NewQueryRedis,
		identity.NewPasswordHasherBcrypt,
		NewConfig,
	)
	return nil, nil, nil
}
