package repoisitory

import (
	"database/sql"
	"log"
)

type DB struct {
	DB *sql.DB
}

var dbConn DB

func ConnectionDB(dsn string) (*DB, error) {
	log.Println("Connecting to Database...")
	d, err := NewMySqlConnection(dsn)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Database!!!")

	dbConn.DB = d
	return &dbConn, nil
}

func NewMySqlConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Println("error")
		return nil, err
	}

	return db, nil
}
