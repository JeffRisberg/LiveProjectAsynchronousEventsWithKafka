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
	var id string

	db, err := sql.Open("mysql", "developer:123456@/liveproject")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select id from processed_events where id=? and event_name=?", event.ID(), event.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	rows.Scan(&id)

	return len(id) > 0, nil
}

// InsertEvent will insert a row into the processed_events table to indicate an event was processed.
func InsertEvent(event events.Event) error {
	db, err := sql.Open("mysql", "developer:123456@/liveproject")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into processed_events (id, event_name, processed_timestamp) values (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(event.ID(), event.Name(), time.Now())
	if err != nil {
		log.Fatal(err)
	}

	return err
}
