package repository

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
)

type BookPageRepository interface {
	CreateBookPage(ctx context.Context, bookPage *entity.BookPage) (string, error)
	GetBookPageByID(ctx context.Context, id string) (entity.BookPage, error)
	GetBookPagesByBookID(ctx context.Context, bookID string) ([]entity.BookPage, error)
	GetBookpagesByBookIDWithPagination(ctx context.Context, bookID string, pageSize int64, pageNumber int64) ([]entity.BookPage, error)
	GetBookpagesByPageRange(ctx context.Context, bookID string, startPage int64, endPage int64) ([]entity.BookPage, error)
	GetBookpagesAroundPageNumber(ctx context.Context, bookID string, centerPage int64, offset int64) ([]entity.BookPage, error)
}
