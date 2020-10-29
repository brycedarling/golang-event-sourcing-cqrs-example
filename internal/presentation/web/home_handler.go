package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/brycedarling/go-practical-microservices/internal/domain/viewing"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
)

// HomeHandler ...
type HomeHandler struct {
	viewingQuery viewing.Query
}

var _ http.Handler = (*HomeHandler)(nil)

// NewHomeHandler ...
func NewHomeHandler(conf *config.Config) *HomeHandler {
	return &HomeHandler{conf.ViewingQuery}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleHome(w, r)
	default:
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HomeHandler) handleHome(w http.ResponseWriter, r *http.Request) {
	v, err := h.viewingQuery.Find()
	if err != nil {
		log.Printf("Unexpected error finding viewing query: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(v)
}
