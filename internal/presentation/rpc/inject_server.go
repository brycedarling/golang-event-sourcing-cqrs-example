//+build wireinject

package rpc

import (
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
	"github.com/google/wire"
)

// InitializeServer ...
func InitializeServer(conf *config.Config) (*Server, error) {
	wire.Build(
		NewServer,
	)
	return nil, nil
}
