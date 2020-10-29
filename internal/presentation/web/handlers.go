package web

import (
	"net/http"
)

// Handlers ...
type Handlers map[string]http.Handler

// NewHandlers ...
func NewHandlers(
	home *HomeHandler,
	viewing *ViewingHandler,
	register *RegisterHandler,
	authentication *AuthenticationHandler,
) Handlers {
	return Handlers{
		"":         home,
		"viewing":  viewing,
		"register": register,
		"login":    authentication,
	}
}
