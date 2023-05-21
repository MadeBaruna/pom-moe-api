package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var pool *pgxpool.Pool

func Connect() {
	p, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
		os.Exit(1)
	}

	pool = p
}

func Close() {
	pool.Close()
}

func Pool() *pgxpool.Pool {
	return pool
}
