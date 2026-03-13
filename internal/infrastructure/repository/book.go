package repository

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
)

type BookRepository interface {
	CreateBook(ctx context.Context, book *entity.Book) (string, error)
	GetBookByID(ctx context.Context, id string) (entity.Book, error)
	GetAllBooks(ctx context.Context) ([]entity.Book, error)
	GetBooksWithPagination(ctx context.Context, pageSize int64, pageNumber int64) ([]entity.Book, error)
}
