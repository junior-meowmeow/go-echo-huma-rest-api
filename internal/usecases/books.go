package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories/mongo_repositories"
)

type BooksUseCase interface {
	CreateBook(ctx context.Context, book *entities.Book) (string, error)
	GetAllBooks(ctx context.Context) ([]entities.Book, error)
	GetBooksWithPagination(ctx context.Context, pageSize int64, pageNumber int64) ([]entities.Book, error)
	GetBookByID(ctx context.Context, id string) (entities.Book, error)
}

type booksUseCase struct {
	BooksRepository mongo_repositories.BooksRepository
}

func NewBooksUseCase(booksRepo mongo_repositories.BooksRepository) *booksUseCase {
	return &booksUseCase{
		BooksRepository: booksRepo,
	}
}

func (u *booksUseCase) CreateBook(ctx context.Context, book *entities.Book) (string, error) {
	currentTime := time.Now()
	book.CreatedAt = currentTime
	book.ModifiedAt = currentTime

	id, err := u.BooksRepository.CreateBook(ctx, book)
	if err != nil {
		return "", fmt.Errorf("failed to create book: %w", err)
	}

	return id, nil
}

func (u *booksUseCase) GetAllBooks(ctx context.Context) ([]entities.Book, error) {
	books, err := u.BooksRepository.GetAllBooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books: %w", err)
	}

	return books, nil
}

func (u *booksUseCase) GetBooksWithPagination(ctx context.Context, pageSize int64, pageNumber int64) ([]entities.Book, error) {
	books, err := u.BooksRepository.GetBooksWithPagination(ctx, pageSize, pageNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books: %w", err)
	}

	return books, nil
}

func (u *booksUseCase) GetBookByID(ctx context.Context, id string) (entities.Book, error) {
	book, err := u.BooksRepository.GetBookByID(ctx, id)
	if err != nil {
		return entities.Book{}, fmt.Errorf("failed to fetch book: %w", err)
	}

	return book, nil
}
