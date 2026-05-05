package integration_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/config"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/api"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/external"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/external/petstore/client"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/postgres"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/storage"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/utility"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/test/helper/adaptor/database"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/test/helper/testenv"
)

type IntegrationTestSuite struct {
	suite.Suite

	Database      database.Adapter
	MongoDB       *mongo.Database
	PgxPool       *pgxpool.Pool
	S3Client      *s3.Client
	MockPetServer *httptest.Server

	Repositories     *repository.Repositories
	Storages         *storage.Storages
	ExternalServices *external.ExternalServices
	Utilities        *utility.Utilities
	Router           http.Handler
}

func (s *IntegrationTestSuite) SetupMongo() {
	s.MongoDB = testenv.SetupMongoDatabase(s.T())
	s.Database = database.NewMongoAdapter(s.MongoDB)

	repos := repository.Repositories{
		Review:     mongodb.NewReviewRepository(s.MongoDB),
		FileRecord: mongodb.NewFileRecordRepository(s.MongoDB),
		User:       mongodb.NewUserRepository(s.MongoDB),
		Book:       mongodb.NewBookRepository(s.MongoDB),
		BookPage:   mongodb.NewBookPageRepository(s.MongoDB),
	}
	s.setupApp(&repos)
}

func (s *IntegrationTestSuite) SetupPostgres() {
	s.PgxPool = testenv.SetupPostgresDatabase(s.T())
	s.Database = database.NewPostgresAdapter(s.PgxPool)

	repos := repository.Repositories{
		Book:     postgres.NewBookRepository(s.PgxPool),
		BookPage: postgres.NewBookPageRepository(s.PgxPool),
	}
	s.setupApp(&repos)
}

func (s *IntegrationTestSuite) setupApp(repos *repository.Repositories) {
	s.Repositories = repos
	s.S3Client = testenv.SetupS3Client(s.T())

	s.MockPetServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}))
	petClient, _ := client.NewClientWithResponses(s.MockPetServer.URL)

	s.Storages = storage.NewStorages(s.S3Client, "test-bucket")
	s.ExternalServices = external.NewExternalServices(petClient)
	s.Utilities = utility.NewUtilities("test-secret", 72*time.Hour)

	usecases := usecase.NewUseCases(s.Repositories, s.Storages, s.ExternalServices, s.Utilities)
	handlers := handler.NewHandlers(usecases)

	appConfig := config.AppConfig{RequestPerSecLimit: 0,
		EnablePrometheus: false}
	s.Router = api.NewRouter(handlers, s.Utilities, appConfig)
}
