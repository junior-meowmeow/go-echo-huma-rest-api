package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func newMongoDBClient(ctx context.Context, user string, pass string, host string, port string) (*mongo.Client, error) {
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port)

	opts := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}
	log.Printf("Created a new MongoDB client and connected to %s:%s\n", host, port)

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
