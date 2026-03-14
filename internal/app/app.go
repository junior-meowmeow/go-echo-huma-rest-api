package app

import (
	"context"
	"log"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/config"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/api"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/external"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/storage"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"

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

	// Initialize External Service Clients
	petStoreClient, err := newPetStoreClient(cfg.PetStoreURL, 5*time.Second)
	if err != nil {
		log.Printf("Failed to initialize PetStore client: %v\n", err)
		return nil, err
	}

	// Initialize Infrastructures
	repositories := repository.NewRepositories(mongoDB)
	storages := storage.NewStorages(s3Client, cfg.S3Bucket)
	externalServices := external.NewExternalServices(petStoreClient)

	// Initialize Use Cases
	usecases := usecase.NewUseCases(repositories, storages, externalServices)

	// Initialize REST API Handlers
	handlers := handler.NewHandlers(usecases)

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
