package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:amoory2003amoory@db:5432/addition?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database: %v\n", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Unable to ping database: %v\n", err)
	}
	DB = pool

	log.Println("Connected to database, successfully")
}

func Close() {
	if DB != nil {
		DB.Close()
		log.Println("Closed database connection")
	}
}
