package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func InitDb(context context.Context) (*Queries, *pgxpool.Pool, error) {
	err := godotenv.Load()
	if err != nil {

	}
	dsn := os.Getenv("DATABASE_URL")
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, nil, err
	}

	pool, err := pgxpool.NewWithConfig(context, config)
	if err != nil {
		return nil, nil, err
	}

	queries := New(pool)
	return queries, pool, nil
}
