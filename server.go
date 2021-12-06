package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

// Server represents the web server hosting the service
type Server struct {
	Port int
}

// Health returns a HTTP 200 status code indicating the service is alive
func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Root returns a HTTP 200 status code
func Root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// ListenAndServe will start the web server and listen for requests
func (s *Server) ListenAndServe() error {

	// setup CHI router
	r := chi.NewRouter()

	// setup middlewares
	r.Use(middleware.Heartbeat("/ping")) // allows LB to verify service up
	r.Use(middleware.RequestID)          // ensures a request ID is logged
	//r.Use(logger.NewStructuredLogger())  // uses structured logging like our app (logs only at debug level)
	r.Use(middleware.Recoverer) // handles any unhandles errors and returns a 500

	// setup supported routes
	r.Get("/", Root)
	r.Get("/health", Health)

	address := fmt.Sprintf(":%d", s.Port)
	log.WithField("address", address).Info("server starting")

	return http.ListenAndServe(address, r)
}
