package database

import (
	"database/sql"
	"log"

	db "osvauld/db/sqlc"

	_ "github.com/lib/pq" // The Postgres driver
)

var (
	DB  *sql.DB
	Q   *db.Queries
	err error
)

func DbConnection(masterDSN string) error {
	// Connecting to master database
	DB, err = sql.Open("postgres", masterDSN)
	if err != nil {
		log.Fatalf("Failed to connect to master database: %v", err)
		return err
	}

	Q = db.New(DB)

	// TODO: Add logic for connecting to replica, if needed

	return nil
}
