package pgsql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var Pool *pgxpool.Pool

func Init() *pgxpool.Pool {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("could not load env vars: %v", err)
	}
	connectionString := os.Getenv("POSTGRES_CONN_STRING")
	var err error
	Pool, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Fatalf("unable to establish PostgreSQL connection: %v", err)
	}
	return Pool
}
