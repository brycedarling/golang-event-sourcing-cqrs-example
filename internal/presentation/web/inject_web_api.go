//+build wireinject

package web

import (
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
	"github.com/google/wire"
)

// InitializeAPI ...
func InitializeAPI(conf *config.Config) (API, func(), error) {
	wire.Build(
		NewHomeHandler,
		NewViewingHandler,
		NewRegisterHandler,
		NewAuthenticationHandler,
		NewHandlers,
		NewListener,
		NewAPI,
	)
	return nil, nil, nil
}
