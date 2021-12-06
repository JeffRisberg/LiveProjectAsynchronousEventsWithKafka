package server

import (
	"fmt"
	"github.com/JeffRisberg/LiveProjectAsynchronousEventsWithKafka/config"
	"github.com/JeffRisberg/LiveProjectAsynchronousEventsWithKafka/events"
	"github.com/JeffRisberg/LiveProjectAsynchronousEventsWithKafka/models"
	"github.com/JeffRisberg/LiveProjectAsynchronousEventsWithKafka/publisher"
	"github.com/google/uuid"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

// Server represents the web server hosting the service
type Server struct {
	Port int
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("The order service is running"))
}

func orders(w http.ResponseWriter, r *http.Request) {
	topic := config.OrderReceivedTopicName

	var order = models.Order{
		ID: uuid.New(),
	}

	var event = events.OrderReceived{
		EventBase: events.BaseEvent{
			EventID:        uuid.New(),
			EventTimestamp: time.Now(),
		},
		EventBody: order,
	}

	publisher.PublishEvent(event, topic)

	w.Write([]byte("Message sent to topic"))
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
	r.Get("/", orders)
	r.Get("/health", health)

	address := fmt.Sprintf(":%d", s.Port)
	log.WithField("address", address).Info("server starting")

	return http.ListenAndServe(address, r)
}
