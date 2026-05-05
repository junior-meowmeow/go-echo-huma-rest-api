package testenv

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // Required for goose to use pgx
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func SetupPostgresDatabase(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	dbName := "test_db"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err, "failed to start postgres container")

	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	connString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "failed to get connection string")

	sqlDB, err := sql.Open("pgx", connString)
	require.NoError(t, err, "failed to open sql db for migrations")
	defer sqlDB.Close()

	rootDir := getProjectRoot()
	migrationDir := filepath.Join(rootDir, "db", "migrations")
	err = goose.Up(sqlDB, migrationDir)
	require.NoError(t, err, "failed to run goose migrations")

	pool, err := pgxpool.New(ctx, connString)
	require.NoError(t, err, "failed to connect pgxpool")

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func CleanPostgresTable(t *testing.T, pool *pgxpool.Pool, tables ...string) {
	t.Helper()
	for _, table := range tables {
		_, err := pool.Exec(context.Background(), "TRUNCATE TABLE "+table+" CASCADE;")
		require.NoError(t, err, "failed to truncate table: "+table)
	}
}
