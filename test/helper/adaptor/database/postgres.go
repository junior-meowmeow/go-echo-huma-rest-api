package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAdapter struct {
	pool *pgxpool.Pool
}

func NewPostgresAdapter(pool *pgxpool.Pool) *PostgresAdapter {
	return &PostgresAdapter{pool: pool}
}
