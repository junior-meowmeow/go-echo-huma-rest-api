package repository

import (
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/port"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/sqlc"
)

type Repositories struct {
	Review     port.ReviewRepository
	FileRecord port.FileRecordRepository
	Book       port.BookRepository
	BookPage   port.BookPageRepository
	User       port.UserRepository
}

func NewRepositories(mongoDB *mongo.Database, pgxPool *pgxpool.Pool) *Repositories {
	return &Repositories{
		Review:     mongodb.NewReviewRepository(mongoDB),
		FileRecord: mongodb.NewFileRecordRepository(mongoDB),
		User:       mongodb.NewUserRepository(mongoDB),
		Book:       sqlc.NewBookRepository(pgxPool),
		BookPage:   sqlc.NewBookPageRepository(pgxPool),
	}
}
