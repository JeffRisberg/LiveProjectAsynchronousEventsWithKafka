package main

import (
	"github.com/JeffRisberg/LiveProjectAsynchronousEventsWithKafka/config"
	"github.com/JeffRisberg/LiveProjectAsynchronousEventsWithKafka/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("startup")

	s := server.Server{
		Port: config.Port(),
	}

	log.Fatal(s.ListenAndServe())
}
