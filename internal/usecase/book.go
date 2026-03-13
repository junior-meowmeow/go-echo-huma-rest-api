package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository"
)

type BookUseCase interface {
	CreateBook(ctx context.Context, book *entity.Book) (string, error)
	GetAllBooks(ctx context.Context) ([]entity.Book, error)
	GetBooksWithPagination(ctx context.Context, pageSize int64, pageNumber int64) ([]entity.Book, error)
	GetBookByID(ctx context.Context, id string) (entity.Book, error)
}

type bookUseCase struct {
	BookRepository repository.BookRepository
}

func NewBookUseCase(bookRepository repository.BookRepository) *bookUseCase {
	return &bookUseCase{
		BookRepository: bookRepository,
	}
}

func (u *bookUseCase) CreateBook(ctx context.Context, book *entity.Book) (string, error) {
	currentTime := time.Now()
	book.CreatedAt = currentTime
	book.ModifiedAt = currentTime

	id, err := u.BookRepository.CreateBook(ctx, book)
	if err != nil {
		return "", fmt.Errorf("failed to create book: %w", err)
	}

	return id, nil
}

func (u *bookUseCase) GetAllBooks(ctx context.Context) ([]entity.Book, error) {
	books, err := u.BookRepository.GetAllBooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books: %w", err)
	}

	return books, nil
}

func (u *bookUseCase) GetBooksWithPagination(ctx context.Context, pageSize int64, pageNumber int64) ([]entity.Book, error) {
	books, err := u.BookRepository.GetBooksWithPagination(ctx, pageSize, pageNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books: %w", err)
	}

	return books, nil
}

func (u *bookUseCase) GetBookByID(ctx context.Context, id string) (entity.Book, error) {
	book, err := u.BookRepository.GetBookByID(ctx, id)
	if err != nil {
		return entity.Book{}, fmt.Errorf("failed to fetch book: %w", err)
	}

	return book, nil
}
