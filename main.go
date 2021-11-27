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
}

func main() {
	log.Info("startup")

	http.HandleFunc("/orders", orders)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
