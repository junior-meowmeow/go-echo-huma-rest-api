package repository

import (
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repositories struct {
	Review     ReviewRepository
	FileRecord FileRecordRepository
	Book       BookRepository
	BookPage   BookPageRepository
}

func NewRepositories(mongoDB *mongo.Database) *Repositories {
	return &Repositories{
		Review:     mongodb.NewReviewRepository(mongoDB),
		FileRecord: mongodb.NewFileRecordRepository(mongoDB),
		Book:       mongodb.NewBookRepository(mongoDB),
		BookPage:   mongodb.NewBookPageRepository(mongoDB),
	}
}
