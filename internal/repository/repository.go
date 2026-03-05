package repository

import (
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository/mongo_repository"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository/s3_repository"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repositories struct {
	Reviews       mongo_repository.ReviewsRepository
	Files         mongo_repository.FileMetadataRepository
	Books         mongo_repository.BooksRepository
	BookPages     mongo_repository.BookPagesRepository
	ObjectStorage s3_repository.ObjectStorage
}

func NewRepositories(mongoDB *mongo.Database, s3Client *s3.Client, bucketName string) *Repositories {
	return &Repositories{
		Reviews:       mongo_repository.NewMongoReviewsRepository(mongoDB),
		Files:         mongo_repository.NewMongoFilesRepository(mongoDB),
		Books:         mongo_repository.NewMongoBooksRepository(mongoDB),
		BookPages:     mongo_repository.NewMongoBookPagesRepository(mongoDB),
		ObjectStorage: s3_repository.NewS3Repository(s3Client, bucketName),
	}
}
