package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type middleware func(http.Handler) http.Handler

func applyMiddleware(h http.Handler, m ...middleware) http.Handler {
	var middleware middleware
	for i := len(m) - 1; i >= 0; i-- {
		middleware = m[i]
		h = middleware(h)
	}
	return h
}

type contextKey string

const contextTraceIDKey contextKey = "traceID"

func addContextTraceID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.NewUUID()
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.WithValue(r.Context(), contextTraceIDKey, id.String())
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func contextTraceID(r *http.Request) string {
	traceID := r.Context().Value(contextTraceIDKey)
	if traceID != nil {
		return traceID.(string)
	}
	return ""
}

const contextUserIDKey contextKey = "userID"

func addContextUserID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID *string = nil

		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			h.ServeHTTP(w, r)
		} else {
			claims, err := ParseJWT(authorization)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
			} else if claims != nil {
				userID = &claims.UserID
				ctx := context.WithValue(r.Context(), contextUserIDKey, userID)
				h.ServeHTTP(w, r.WithContext(ctx))
			}
		}
	})
}

func contextUserID(r *http.Request) *string {
	userID := r.Context().Value(contextUserIDKey)
	if userID != nil {
		return userID.(*string)
	}
	return nil
}

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := contextTraceID(r)
		start := time.Now()
		url := r.URL.String()
		log.Printf("-> %s %s | UserID %v | TraceID %s", r.Method, url, contextUserID(r), traceID)
		lrw := newLoggingResponseWriter(w)
		h.ServeHTTP(lrw, r)
		duration := time.Now().Sub(start)
		log.Printf("<- %s %s responded with %s, took %s | TraceID %s",
			r.Method, url, http.StatusText(lrw.statusCode), duration, traceID)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
