package app

import (
	"context"
	"log"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/config"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controllers/restapi/api"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controllers/restapi/handlers"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecases"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Application struct {
	Router      *echo.Echo
	mongoClient *mongo.Client
}

func NewApplication(ctx context.Context, cfg config.Config) (*Application, error) {
	// Initialize MongoDB
	mongoClient, err := newMongoDBClient(ctx, cfg.MongoUser, cfg.MongoPass, cfg.MongoHost, cfg.MongoPort)
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v\n", err)
		return nil, err
	}
	mongoDB := mongoClient.Database(cfg.DBName)

	// Initialize S3
	s3Client, err := newS3Client(ctx, cfg.S3Endpoint)
	if err != nil {
		log.Printf("Failed to connect to S3: %v\n", err)
		return nil, err
	}

	// Initialize Repositories
	repositories := repositories.NewRepositories(mongoDB, s3Client, cfg.S3Bucket)

	// Initialize Use Cases
	usecases := usecases.NewUseCases(repositories)

	// Initialize REST API Handlers
	handlers := handlers.NewHandlers(usecases)

	// Initialize Router and Register APIs
	router := api.NewRouter(handlers, cfg.APIBasePath)

	// Initialize Application
	application := Application{
		Router:      router,
		mongoClient: mongoClient,
	}

	return &application, nil
}

func (a *Application) GracefulShutdown(ctx context.Context) {
	err := disconnectMongoDB(ctx, a.mongoClient)
	if err != nil {
		log.Printf("Error disconnecting MongoDB: %v\n", err)
	}
}
