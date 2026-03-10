package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories/mongo_repositories"
)

type BookPagesUseCase interface {
	CreateBookPage(ctx context.Context, bookPage *entities.BookPage) (string, error)
	GetAllBookPages(ctx context.Context, bookID string) ([]entities.BookPage, error)
	GetBookPagesWithPagination(ctx context.Context, bookID string, pageSize int64, pageNumber int64) ([]entities.BookPage, error)
	GetBookPagesByRange(ctx context.Context, bookID string, startPage int64, endPage int64) ([]entities.BookPage, error)
	GetBookPagesByOffset(ctx context.Context, bookID string, centerPage int64, offset int64) ([]entities.BookPage, error)
	GetBookPageByID(ctx context.Context, id string) (entities.BookPage, error)
}

type bookPagesUseCase struct {
	BooksRepository     mongo_repositories.BooksRepository
	BookPagesRepository mongo_repositories.BookPagesRepository
}

func NewBookPagesUseCase(booksRepo mongo_repositories.BooksRepository, bookPagesRepo mongo_repositories.BookPagesRepository) *bookPagesUseCase {
	return &bookPagesUseCase{
		BooksRepository:     booksRepo,
		BookPagesRepository: bookPagesRepo,
	}
}

func (u *bookPagesUseCase) CreateBookPage(ctx context.Context, bookPage *entities.BookPage) (string, error) {
	_, err := u.BooksRepository.GetBookByID(ctx, bookPage.BookID.Hex())
	if err != nil {
		return "", fmt.Errorf("failed to fetch book info: %w", err)
	}

	currentTime := time.Now()
	bookPage.CreatedAt = currentTime
	bookPage.ModifiedAt = currentTime

	id, err := u.BookPagesRepository.CreateBookPage(ctx, bookPage)
	if err != nil {
		return "", fmt.Errorf("failed to create book page: %w", err)
	}

	return id, nil
}

func (u *bookPagesUseCase) GetAllBookPages(ctx context.Context, bookID string) ([]entities.BookPage, error) {
	bookPages, err := u.BookPagesRepository.GetBookPagesByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages: %w", err)
	}

	return bookPages, nil
}

func (u *bookPagesUseCase) GetBookPagesWithPagination(ctx context.Context, bookID string, pageSize int64, pageNumber int64) ([]entities.BookPage, error) {
	bookPages, err := u.BookPagesRepository.GetBookpagesByBookIDWithPagination(ctx, bookID, pageSize, pageNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages: %w", err)
	}

	return bookPages, nil
}

func (u *bookPagesUseCase) GetBookPagesByRange(ctx context.Context, bookID string, startPage int64, endPage int64) ([]entities.BookPage, error) {
	bookPages, err := u.BookPagesRepository.GetBookpagesByPageRange(ctx, bookID, startPage, endPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages by range: %w", err)
	}

	return bookPages, nil
}

func (u *bookPagesUseCase) GetBookPagesByOffset(ctx context.Context, bookID string, centerPage int64, offset int64) ([]entities.BookPage, error) {
	bookPages, err := u.BookPagesRepository.GetBookpagesAroundPageNumber(ctx, bookID, centerPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages by offset: %w", err)
	}

	return bookPages, nil
}

func (u *bookPagesUseCase) GetBookPageByID(ctx context.Context, id string) (entities.BookPage, error) {
	bookPage, err := u.BookPagesRepository.GetBookPageByID(ctx, id)
	if err != nil {
		return entities.BookPage{}, fmt.Errorf("failed to fetch book page: %w", err)
	}

	return *bookPage, nil
}
