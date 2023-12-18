package database

import (
	"database/sql"
	"log"

	db "osvauld/db/sqlc"

	_ "github.com/lib/pq" // The Postgres driver
)

var (
	Store db.Store
)

func DbConnection(masterDSN string) error {
	// Connecting to master database

	conn, err := sql.Open("postgres", masterDSN)
	if err != nil {
		log.Fatalf("Failed to connect to master database: %v", err)
		return err
	}

	if err := conn.Ping(); err != nil {
		log.Fatalf("Failed to connect to master database: %v", err)
		return err
	}

	Store = db.NewStore(conn)

	// TODO: Add logic for connecting to replica, if needed

	return nil
}
