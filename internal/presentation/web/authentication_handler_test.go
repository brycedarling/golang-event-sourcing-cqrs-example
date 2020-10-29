package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
)

var conf *config.Config

func init() {
	conf = config.InitializeTestEnvConfig()
}

func TestAuthenticationHandler(t *testing.T) {
	r, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler := NewAuthenticationHandler(conf)
	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
	}
}
