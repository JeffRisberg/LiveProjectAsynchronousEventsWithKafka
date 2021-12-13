package server

import (
	"encoding/json"
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

// Root returns a HTTP 200 status code
func Root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Health returns a HTTP 200 status code indicating the service is alive
func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// HandleError will publish an error event to Kafka
func HandleError(event events.Event) {
	var err error

	e := translateToErrorEvent(event)
	if err = publisher.PublishEvent(e, config.ErrorsTopicName); err != nil {
		log.WithField("error", err).
			WithField("topic", config.ErrorsTopicName).
			Error("an issue ocurred publishing an error event to Kafka")
	}
}

func translateToErrorEvent(event events.Event) events.Event {
	return events.Error{
		EventBase: events.BaseEvent{
			EventID:        uuid.New(),
			EventTimestamp: time.Now(),
		},
		EventBody: event,
	}
}

// ReceiveOrder handler will accept an order, validate the payload and publish an OrderReceived event to Kafka.
// returns a HTTP 201 status code indicating an order was created
func ReceiveOrder(w http.ResponseWriter, r *http.Request) {
	var o models.Order
	// Create a new ID for the order
	o.ID = uuid.New()

	var err error

	if err = json.NewDecoder(r.Body).Decode(&o); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	log.WithField("order", o).Info("received new order")

	if err = validate(o); err != nil {
		log.WithField("orderID", o.ID).Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)

		// publish to deadletterqueue
		return
	}

	e := translateOrderToReceivedEvent(o)

	log.WithField("event", e).Info("transformed order to event")

	if err = publisher.PublishEvent(e, config.OrderReceivedTopicName); err != nil {
		log.WithField("orderID", o.ID).Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	log.WithField("event", e).Info("published event")

	e = translateOrderToConfirmedEvent(o)

	// publish to order Confirmed
	if err = publisher.PublishEvent(e, config.OrderConfirmedTopicName); err != nil {
		log.WithField("orderID", o.ID).Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Validates the order payload has the necessary information and returns an error if it is invalid
func validate(o models.Order) error {
	if len(o.Products) == 0 {
		return fmt.Errorf("there are no products in the order")
	}

	for i, p := range o.Products {
		if len(p.ProductCode) == 0 {
			return fmt.Errorf("product code is required for product [%d]", i)
		}

		if p.Quantity <= 0 {
			return fmt.Errorf("quantity should be greater than zero for product [%s]", p.ProductCode)
		}
	}

	if len(o.Customer.EmailAddress) == 0 {
		return fmt.Errorf("email address is required")
	}

	if len(o.Customer.ShippingAddress.Line1) == 0 {
		return fmt.Errorf("shipping address line 1 is required")
	}

	if len(o.Customer.ShippingAddress.City) == 0 {
		return fmt.Errorf("shipping address city is required")
	}

	if len(o.Customer.ShippingAddress.PostalCode) == 0 {
		return fmt.Errorf("shipping address postal code is required")
	}

	return nil
}

func translateOrderToReceivedEvent(o models.Order) events.Event {
	var event = events.OrderReceived{
		EventBase: events.BaseEvent{
			EventID:        uuid.New(),
			EventTimestamp: time.Now(),
		},
		EventBody: o,
	}

	return event
}

func translateOrderToConfirmedEvent(o models.Order) events.Event {
	var event = events.OrderConfirmed{
		EventBase: events.BaseEvent{
			EventID:        uuid.New(),
			EventTimestamp: time.Now(),
		},
		EventBody: o,
	}

	return event
}

// ListenAndServe starts the server and listen for requests
func (s *Server) ListenAndServe() error {

	// set up CHI router
	r := chi.NewRouter()

	// set up middlewares
	r.Use(middleware.Heartbeat("/ping")) // allows LB to verify service up
	r.Use(middleware.RequestID)          // ensures a request ID is logged
	//r.Use(logger.NewStructuredLogger())  // uses structured logging like our app (logs only at debug level)
	r.Use(middleware.Recoverer) // handles any unhandles errors and returns a 500

	// set up supported routes
	r.Get("/", Root)
	r.Get("/health", Health)
	r.Post("/orders", ReceiveOrder)

	address := fmt.Sprintf(":%d", s.Port)
	log.WithField("address", address).Info("server starting")

	return http.ListenAndServe(address, r)
}
