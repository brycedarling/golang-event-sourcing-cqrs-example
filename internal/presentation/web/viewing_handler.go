package web

import (
	"log"
	"net/http"

	"github.com/brycedarling/go-practical-microservices/internal/domain/viewing/command"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
)

// ViewingHandler ...
type ViewingHandler struct {
	eventStore eventstore.Store
}

var _ http.Handler = (*ViewingHandler)(nil)

// NewViewingHandler ...
func NewViewingHandler(conf *config.Config) *ViewingHandler {
	return &ViewingHandler{conf.EventStore}
}

func (h *ViewingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var id, junk string
	id, r.URL.Path = shiftPath(r.URL.Path)
	junk, r.URL.Path = shiftPath(r.URL.Path)
	if id == "" || junk != "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case "POST":
		h.handleViewVideo(w, r, id)
	default:
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ViewingHandler) handleViewVideo(w http.ResponseWriter, r *http.Request, id string) {
	cmd, err := command.NewViewVideoCommand(h.eventStore, contextTraceID(r), contextUserID(r), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = cmd.Execute()
	if err != nil {
		log.Printf("Error viewing video: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
