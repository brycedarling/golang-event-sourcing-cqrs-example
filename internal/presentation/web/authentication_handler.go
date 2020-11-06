package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/brycedarling/go-practical-microservices/internal/domain/identity"
	"github.com/brycedarling/go-practical-microservices/internal/domain/identity/command"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
)

// AuthenticationHandler ...
type AuthenticationHandler struct {
	eventStore     eventstore.Store
	identityQuery  identity.Query
	passwordHasher identity.PasswordHasher
}

var _ http.Handler = (*AuthenticationHandler)(nil)

// NewAuthenticationHandler ...
func NewAuthenticationHandler(conf *config.Config) *AuthenticationHandler {
	return &AuthenticationHandler{conf.EventStore, conf.IdentityQuery, conf.PasswordHasher}
}

func (h *AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var junk string
	junk, r.URL.Path = shiftPath(r.URL.Path)
	if junk != "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case "POST":
		h.handleAuthentication(w, r)
	default:
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AuthenticationHandler) handleAuthentication(w http.ResponseWriter, r *http.Request) {
	cmd, err := h.authenticateCommand(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := cmd.Execute()
	if err != nil {
		if _, ok := err.(command.ErrAuthenticationFailed); ok {
			http.Error(w, "", http.StatusForbidden)
		} else {
			log.Println("Unexpected error authenticating:", err)
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	if err := h.writeJWT(w, id); err != nil {
		log.Println("Unexpected error writing JWT:", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func (h *AuthenticationHandler) authenticateCommand(r *http.Request) (command.AuthenticateCommand, error) {
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

	traceID := contextTraceID(r)

	return command.NewAuthenticateCommand(h.eventStore, h.identityQuery, h.passwordHasher,
		traceID, *input.Email, *input.Password)
}

type authenticateInput struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (h *AuthenticationHandler) decodeInput(r *http.Request) (input authenticateInput, err error) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	err = d.Decode(&input)
	if d.More() {
		err = errors.New("extraneous data after JSON object")
	}
	return
}

func (h *AuthenticationHandler) writeJWT(w http.ResponseWriter, id *identity.Identity) error {
	signedToken, err := SignJWT(id.UserID)
	if err != nil {
		return err
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", signedToken))
	w.WriteHeader(http.StatusOK)
	return nil
}
