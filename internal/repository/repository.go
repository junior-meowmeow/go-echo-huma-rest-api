package repository

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repositories struct {
	Reviews       ReviewsRepository
	Files         FileMetadataRepository
	Books         BooksRepository
	BookPages     BookPagesRepository
	ObjectStorage ObjectStorage
}

func NewRepositories(mongoDB *mongo.Database, s3Client *s3.Client, bucketName string) *Repositories {
	return &Repositories{
		Reviews:       NewMongoReviewsRepository(mongoDB),
		Files:         NewMongoFilesRepository(mongoDB),
		Books:         NewMongoBooksRepository(mongoDB),
		BookPages:     NewMongoBookPagesRepository(mongoDB),
		ObjectStorage: NewS3Repository(s3Client, bucketName),
	}
}
