package web

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
)

// API ...
type API interface {
	Listen()
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// NewAPI ...
func NewAPI(conf *config.Config, h Handlers) (API, func(), error) {
	l, shutdownListener, err := NewTCPListener(fmt.Sprintf(":%s", conf.Env.Port))
	if err != nil {
		return nil, nil, err
	}

	return &api{
		env:      conf.Env.Env,
		handlers: h,
		listener: l,
	}, shutdownListener, nil
}

type api struct {
	env      string
	handlers Handlers
	listener net.Listener
}

var _ API = (*api)(nil)

// Listen ...
func (api *api) Listen() {
	log.Printf("Starting web API in %s on %s", api.env, api.listener.Addr())
	err := http.Serve(api.listener, withGlobalMiddleware(api))
	if err != nil && err != ErrShutdown {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (api *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	if handler, ok := api.handlers[head]; ok {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}
}

func withGlobalMiddleware(api *api) http.Handler {
	return applyMiddleware(api, addContextTraceID, logRequest, addContextUserID)
}
