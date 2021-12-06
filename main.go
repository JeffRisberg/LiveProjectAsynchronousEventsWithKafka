package main

import (
	"github.com/JeffRisberg/LiveProjectAsynchronousEventsWithKafka/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("startup")

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
