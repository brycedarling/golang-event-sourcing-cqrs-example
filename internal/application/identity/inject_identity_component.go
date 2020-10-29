//+build wireinject

package identity

import (
	"github.com/brycedarling/go-practical-microservices/internal/domain/identity"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
	"github.com/google/wire"
)

// InitializeComponent ...
func InitializeComponent(conf *config.Config) *Component {
	wire.Build(
		identity.NewProjector,
		NewRegisterEventHandler,
		NewComponent,
	)
	return nil
}
