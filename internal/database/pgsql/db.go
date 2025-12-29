package pgsql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var Pool *pgxpool.Pool

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConnection(config *Config) *pgxpool.Pool {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("could not load env vars: %v", err)
	}

	// For production
	//connStr := NewConnStr(config)
	connStr := os.Getenv("POSTGRES_CONN_STRING")
	var err error
	Pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("unable to establish PostgreSQL connection: %v", err)
	}
	return Pool
}

func NewConnStr(config *Config) string {
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@localhost:%d/%s",
		config.User,
		config.Password,
		config.Port,
		config.DBName,
	)
	return connStr
}
