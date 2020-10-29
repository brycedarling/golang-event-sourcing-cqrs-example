//+build wireinject

package viewing

import (
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
	"github.com/google/wire"
)

// InitializeAggregator ...
func InitializeAggregator(conf *config.Config) *Aggregator {
	wire.Build(NewAggregator)
	return nil
}
