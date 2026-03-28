package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/config"
)

func newMongoDBClient(ctx context.Context, cfg config.MongoConfig) (*mongo.Client, error) {
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPass, cfg.Host, cfg.Port, cfg.DBName)

	opts := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}
	log.Printf("Created a new MongoDB client and connected to %s:%s\n", cfg.Host, cfg.Port)

	if err := pingMongoDB(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to ping mongoDB: %w", err)
	}

	return client, nil
}

func pingMongoDB(ctx context.Context, client *mongo.Client) error {
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return client.Ping(pingCtx, nil)
}

func disconnectMongoDB(ctx context.Context, client *mongo.Client) error {
	if client == nil {
		log.Println("MongoDB Client is nil.")
		return nil
	}

	err := client.Disconnect(ctx)
	if err != nil {
		return err
	}
	log.Println("MongoDB Client disconnected.")
	return nil
}
