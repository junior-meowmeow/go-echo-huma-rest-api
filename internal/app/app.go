package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/config"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/api"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/external"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/storage"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/utility"
)

type Application struct {
	Router      *echo.Echo
	mongoClient *mongo.Client
}

func NewApplication(ctx context.Context, cfg config.Config) (*Application, error) {
	// Initialize MongoDB
	mongoClient, err := newMongoDBClient(ctx, cfg.Mongo)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MongoDB client: %w", err)
	}
	mongoDB := mongoClient.Database(cfg.Mongo.DBName)

	// Initialize S3
	s3Client, err := newS3Client(ctx, cfg.S3)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize S3 client: %w", err)
	}

	// Initialize External Service Clients
	const clientTimeout = 5 * time.Second
	petStoreClient, err := newPetStoreClient(cfg.Client.PetStoreURL, clientTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PetStore client: %w", err)
	}

	// Initialize Infrastructures
	repositories := repository.NewRepositories(mongoDB)
	storages := storage.NewStorages(s3Client, cfg.S3.Bucket)
	externalServices := external.NewExternalServices(petStoreClient)

	// Initialize Utilities
	utilities := utility.NewUtilities(cfg.Auth.JWTSecret, cfg.Auth.TokenExpiration)

	// Initialize Use Cases
	usecases := usecase.NewUseCases(repositories, storages, externalServices, utilities)

	// Initialize REST API Handlers
	handlers := handler.NewHandlers(usecases)

	// Initialize REST API Router and Register APIs
	router := api.NewRouter(handlers, utilities, cfg.App)

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
		slog.ErrorContext(ctx, "Error disconnecting MongoDB", slog.Any("error", err))
	}
}
