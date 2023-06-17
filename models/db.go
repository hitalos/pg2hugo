package models

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

// Connect starts a connection pool with the database
func Connect() error {
	cfg, err := pgxpool.ParseConfig(os.ExpandEnv(os.Getenv("DSN")))
	if err != nil {
		log.Fatalf("Unable to parse config env: %s", err)
	}

	log.Println("Connecting with database")
	db, err = pgxpool.NewWithConfig(context.Background(), cfg)
	return err
}
