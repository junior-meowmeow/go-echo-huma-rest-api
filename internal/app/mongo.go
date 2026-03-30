package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/config"
)

func newMongoDBClient(ctx context.Context, cfg config.MongoConfig) (*mongo.Client, error) {
	mongoHostPort := net.JoinHostPort(cfg.Host, cfg.Port)
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s/%s", cfg.DBUser, cfg.DBPass, mongoHostPort, cfg.DBName)

	opts := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("Created a new MongoDB client and connected to %s", mongoHostPort))

	if err := pingMongoDB(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to ping mongoDB: %w", err)
	}

	return client, nil
}

func pingMongoDB(ctx context.Context, client *mongo.Client) error {
	const pingTimeout = 5 * time.Second
	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	return client.Ping(pingCtx, nil)
}

func disconnectMongoDB(ctx context.Context, client *mongo.Client) error {
	if client == nil {
		slog.DebugContext(ctx, "MongoDB Client is nil.")
		return nil
	}

	err := client.Disconnect(ctx)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "MongoDB Client disconnected.")
	return nil
}
