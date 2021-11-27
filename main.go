package main

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/JeffRisberg/LiveProjectAsynchronousEventsWithKafka/config"
	"github.com/JeffRisberg/LiveProjectAsynchronousEventsWithKafka/events"
	"github.com/JeffRisberg/LiveProjectAsynchronousEventsWithKafka/models"
	"github.com/JeffRisberg/LiveProjectAsynchronousEventsWithKafka/publisher"
)

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

func main() {
	log.Info("startup")

	server := &http.Server{
		Addr: "0.0.0.0:8080",
	}

	http.HandleFunc("/orders", orders)
	http.HandleFunc("/health", health)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
