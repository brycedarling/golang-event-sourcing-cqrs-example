// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package rpc

import (
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
)

// Injectors from inject_server.go:

// InitializeServer ...
func InitializeServer(conf *config.Config) (*Server, error) {
	server := NewServer(conf)
	return server, nil
}
