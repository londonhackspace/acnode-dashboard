package main

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

type statusIntercepter struct {
	next http.ResponseWriter
	statuscode *int
}

func (h statusIntercepter) Header() http.Header {
	return h.next.Header()
}

func (h statusIntercepter) Write(data []byte) (int, error) {
	return h.next.Write(data)
}

func (h statusIntercepter) WriteHeader(statusCode int) {
	h.next.WriteHeader(statusCode)
	if h.statuscode != nil {
		*h.statuscode = statusCode
	}
}

func createStatusIntercepter(next http.ResponseWriter, i *int) statusIntercepter {
	return statusIntercepter{
		next: next,
		statuscode: i,
	}
}

type LoggingHandler struct {
	next http.Handler
}

func (h LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	intercepter := createStatusIntercepter(w, &statusCode)
	h.next.ServeHTTP(intercepter, r)

	log.Debug().Str("Remote", r.RemoteAddr).
		Str("Method", r.Method).
		Str("UserAgent", r.UserAgent()).
		Str("Path", r.URL.Path).
		Int("StatusCode", statusCode).
		Msg("Serving page")
}
