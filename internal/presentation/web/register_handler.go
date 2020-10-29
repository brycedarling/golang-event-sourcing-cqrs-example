package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/brycedarling/go-practical-microservices/internal/domain/identity"
	domIdentity "github.com/brycedarling/go-practical-microservices/internal/domain/identity"
	"github.com/brycedarling/go-practical-microservices/internal/domain/identity/command"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
)

// RegisterHandler ...
type RegisterHandler struct {
	eventStore     eventstore.Store
	identityQuery  identity.Query
	passwordHasher identity.PasswordHasher
}

var _ http.Handler = (*RegisterHandler)(nil)

// NewRegisterHandler ...
func NewRegisterHandler(conf *config.Config) *RegisterHandler {
	return &RegisterHandler{conf.EventStore, conf.IdentityQuery, conf.PasswordHasher}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var junk string
	junk, r.URL.Path = shiftPath(r.URL.Path)
	if junk != "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case "POST":
		h.handleRegister(w, r)
	default:
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RegisterHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	cmd, err := h.registerCommand(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = cmd.Execute()
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		return
	}
	if err == domIdentity.ErrIdentityAlreadyExists {
		http.Error(w, "", http.StatusUnprocessableEntity)
	} else {
		log.Println("Unexpected error registering user:", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func (h *RegisterHandler) registerCommand(r *http.Request) (command.RegisterCommand, error) {
	input, err := h.decodeInput(r)
	if err != nil {
		return nil, err
	}
	if input.Email == nil {
		return nil, errors.New("missing email")
	}
	if input.Password == nil {
		return nil, errors.New("missing password")
	}

	return command.NewRegisterCommand(h.eventStore, h.identityQuery, h.passwordHasher,
		contextTraceID(r), *input.Email, *input.Password)
}

type registerInput struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (h *RegisterHandler) decodeInput(r *http.Request) (input registerInput, err error) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	err = d.Decode(&input)
	if d.More() {
		err = errors.New("extraneous data after JSON object")
	}
	return
}
