package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository/mongodb"
)

type BookPageUseCase interface {
	CreateBookPage(ctx context.Context, bookPage *entity.BookPage) (string, error)
	GetAllBookPages(ctx context.Context, bookID string) ([]entity.BookPage, error)
	GetBookPagesWithPagination(ctx context.Context, bookID string, pageSize int64, pageNumber int64) ([]entity.BookPage, error)
	GetBookPagesByRange(ctx context.Context, bookID string, startPage int64, endPage int64) ([]entity.BookPage, error)
	GetBookPagesByOffset(ctx context.Context, bookID string, centerPage int64, offset int64) ([]entity.BookPage, error)
	GetBookPageByID(ctx context.Context, id string) (entity.BookPage, error)
}

type bookPageUseCase struct {
	BooksRepository     mongodb.BookRepository
	BookPagesRepository mongodb.BookPageRepository
}

func NewBookPageUseCase(booksRepo mongodb.BookRepository, bookPagesRepo mongodb.BookPageRepository) *bookPageUseCase {
	return &bookPageUseCase{
		BooksRepository:     booksRepo,
		BookPagesRepository: bookPagesRepo,
	}
}

func (u *bookPageUseCase) CreateBookPage(ctx context.Context, bookPage *entity.BookPage) (string, error) {
	_, err := u.BooksRepository.GetBookByID(ctx, bookPage.BookID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch book: %w", err)
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

func (u *bookPageUseCase) GetAllBookPages(ctx context.Context, bookID string) ([]entity.BookPage, error) {
	bookPages, err := u.BookPagesRepository.GetBookPagesByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages: %w", err)
	}

	return bookPages, nil
}

func (u *bookPageUseCase) GetBookPagesWithPagination(ctx context.Context, bookID string, pageSize int64, pageNumber int64) ([]entity.BookPage, error) {
	bookPages, err := u.BookPagesRepository.GetBookpagesByBookIDWithPagination(ctx, bookID, pageSize, pageNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages: %w", err)
	}

	return bookPages, nil
}

func (u *bookPageUseCase) GetBookPagesByRange(ctx context.Context, bookID string, startPage int64, endPage int64) ([]entity.BookPage, error) {
	bookPages, err := u.BookPagesRepository.GetBookpagesByPageRange(ctx, bookID, startPage, endPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages by range: %w", err)
	}

	return bookPages, nil
}

func (u *bookPageUseCase) GetBookPagesByOffset(ctx context.Context, bookID string, centerPage int64, offset int64) ([]entity.BookPage, error) {
	bookPages, err := u.BookPagesRepository.GetBookpagesAroundPageNumber(ctx, bookID, centerPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages by offset: %w", err)
	}

	return bookPages, nil
}

func (u *bookPageUseCase) GetBookPageByID(ctx context.Context, id string) (entity.BookPage, error) {
	bookPage, err := u.BookPagesRepository.GetBookPageByID(ctx, id)
	if err != nil {
		return entity.BookPage{}, fmt.Errorf("failed to fetch book page: %w", err)
	}

	return bookPage, nil
}
