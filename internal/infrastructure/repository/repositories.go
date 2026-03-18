package repository

import (
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/port"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repositories struct {
	Review     port.ReviewRepository
	FileRecord port.FileRecordRepository
	Book       port.BookRepository
	BookPage   port.BookPageRepository
	User       port.UserRepository
}

func NewRepositories(mongoDB *mongo.Database) *Repositories {
	return &Repositories{
		Review:     mongodb.NewReviewRepository(mongoDB),
		FileRecord: mongodb.NewFileRecordRepository(mongoDB),
		Book:       mongodb.NewBookRepository(mongoDB),
		BookPage:   mongodb.NewBookPageRepository(mongoDB),
		User:       mongodb.NewUserRepository(mongoDB),
	}
}
