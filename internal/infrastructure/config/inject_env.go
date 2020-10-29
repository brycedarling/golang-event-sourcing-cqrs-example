//+build wireinject

package config

import (
	"github.com/google/wire"
)

// InitializeEnv ...
func InitializeEnv() (*Env, error) {
	wire.Build(NewEnv)
	return nil, nil
}
