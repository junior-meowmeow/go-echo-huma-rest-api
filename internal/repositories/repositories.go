package repositories

import (
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories/mongo_repositories"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories/s3_repositories"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repositories struct {
	Reviews       mongo_repositories.ReviewsRepository
	Files         mongo_repositories.FileMetadataRepository
	Books         mongo_repositories.BooksRepository
	BookPages     mongo_repositories.BookPagesRepository
	ObjectStorage s3_repositories.ObjectStorage
}

func NewRepositories(mongoDB *mongo.Database, s3Client *s3.Client, bucketName string) *Repositories {
	return &Repositories{
		Reviews:       mongo_repositories.NewMongoReviewsRepository(mongoDB),
		Files:         mongo_repositories.NewMongoFilesRepository(mongoDB),
		Books:         mongo_repositories.NewMongoBooksRepository(mongoDB),
		BookPages:     mongo_repositories.NewMongoBookPagesRepository(mongoDB),
		ObjectStorage: s3_repositories.NewS3Repository(s3Client, bucketName),
	}
}
