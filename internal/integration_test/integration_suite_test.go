package integration_test

import (
	"net/http"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controllers/restapi/api"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controllers/restapi/handlers"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IntegrationTestSuite struct {
	suite.Suite

	MongoDB  *mongo.Database
	S3Client *s3.Client

	Repositories *repositories.Repositories
	Router       http.Handler
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.MongoDB = setupMongoDatabase(s.T())
	s.S3Client = setupS3Client(s.T())

	s.Repositories = repositories.NewRepositories(s.MongoDB, s.S3Client, "test-bucket")

	handlers := handlers.NewHandlers(s.Repositories)
	s.Router = api.NewRouter(handlers, "")
}
