package server

import (
	"database/sql"
	"time"

	"github.com/JeffRisberg/LiveProjectAsynchronousEventsWithKafka/config"
	"github.com/JeffRisberg/LiveProjectAsynchronousEventsWithKafka/events"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

// DB represents an interaction with a database
type DB struct {
	Address  string
	Username string
	Password string
	Database string
}

// NewDB returns a default localhost reference to a DB
func NewDB() DB {
	return DB{
		Address:  config.DatabaseAddress(),
		Username: config.DatabaseUsername(),
		Password: config.DatabasePassword(),
		Database: config.DatabaseName(),
	}
}

// EventExists will check to see if an event has already been processed
func EventExists(event events.Event) (bool, error) {

	db, err := sql.Open("mysql", "developer:123456@/dbname")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select id from events.processed_events where id=$1 and event_name=$2", event.ID(), event.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return false, nil
}

// InsertEvent will insert a row into the processed_events table to indicate an event was processed.
func InsertEvent(event events.Event) error {
	db, err := sql.Open("mysql", "developer:123456@/dbname")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into events.processed_events (id, event_name, processed_timestamp) values ($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(event.ID(), event.Name(), time.Now())
	if err != nil {
		log.Fatal(err)
	}

	return err
}
