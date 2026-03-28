package integration_test

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/suite"
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

type IntegrationTestSuite struct {
	suite.Suite

	MongoDB  *mongo.Database
	S3Client *s3.Client

	Repositories     *repository.Repositories
	Storages         *storage.Storages
	ExternalServices *external.ExternalServices
	Utilities        *utility.Utilities
	Router           http.Handler
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.MongoDB = setupMongoDB(s.T())
	s.S3Client = setupS3Client(s.T())

	s.Repositories = repository.NewRepositories(s.MongoDB)
	s.Storages = storage.NewStorages(s.S3Client, "test-bucket")
	s.ExternalServices = external.NewExternalServices(nil)
	s.Utilities = utility.NewUtilities("test-secret")

	usecases := usecase.NewUseCases(s.Repositories, s.Storages, s.ExternalServices, s.Utilities)
	handlers := handler.NewHandlers(usecases)
	s.Router = api.NewRouter(handlers, s.Utilities, config.AppConfig{})
}
