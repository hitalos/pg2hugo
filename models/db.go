package models

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

var db *pgxpool.Pool

// Connect starts a connection pool with the database
func Connect() error {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error on read '.env': %q", err)
		log.Printf("value of DSN will be: %q", os.Getenv("DSN"))
	}
	log.Println("Connecting with database")
	db, err = pgxpool.Connect(context.Background(), os.Getenv("DSN"))
	return err
}
