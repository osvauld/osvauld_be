package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // The Postgres driver
)

var (
	DB  *sql.DB
	err error
)

func DbConnection(masterDSN string) error {
	// Connecting to master database
	DB, err = sql.Open("postgres", masterDSN)
	if err != nil {
		log.Fatalf("Failed to connect to master database: %v", err)
		return err
	}

	// TODO: Add logic for connecting to replica, if needed

	return nil
}
