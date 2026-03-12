package repository

import (
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository/mongodb"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository/objectstorage"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repositories struct {
	Review        mongodb.ReviewRepository
	FileRecord    mongodb.FileRecordRepository
	Book          mongodb.BookRepository
	BookPage      mongodb.BookPageRepository
	ObjectStorage objectstorage.ObjectStorage
}

func NewRepositories(mongoDB *mongo.Database, s3Client *s3.Client, bucketName string) *Repositories {
	return &Repositories{
		Review:        mongodb.NewReviewRepository(mongoDB),
		FileRecord:    mongodb.NewFileRecordRepository(mongoDB),
		Book:          mongodb.NewBookRepository(mongoDB),
		BookPage:      mongodb.NewBookPagesRepository(mongoDB),
		ObjectStorage: objectstorage.NewS3Repository(s3Client, bucketName),
	}
}
