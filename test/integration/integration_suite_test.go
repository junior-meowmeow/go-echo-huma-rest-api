package integration_test

import (
	"net/http"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/api"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/storage"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IntegrationTestSuite struct {
	suite.Suite

	MongoDB  *mongo.Database
	S3Client *s3.Client

	Repositories *repository.Repositories
	Storages     *storage.Storages
	Router       http.Handler
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.MongoDB = setupMongoDB(s.T())
	s.S3Client = setupS3Client(s.T())

	s.Repositories = repository.NewRepositories(s.MongoDB)
	s.Storages = storage.NewStorages(s.S3Client, "test-bucket")

	usecases := usecase.NewUseCases(s.Repositories, s.Storages)
	handlers := handler.NewHandlers(usecases)
	s.Router = api.NewRouter(handlers, "")
}
