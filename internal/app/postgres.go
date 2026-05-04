package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/config"
)

func newPostgresClient(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	hostPort := net.JoinHostPort(cfg.Host, cfg.Port)
	postgresDSN := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPass, hostPort, cfg.DBName)

	poolConfig, err := pgxpool.ParseConfig(postgresDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("Created a new PostgreSQL client and connected to %s", hostPort))

	return pool, nil
}
