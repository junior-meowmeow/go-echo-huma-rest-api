package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/port"
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
	BookRepository     port.BookRepository
	BookPageRepository port.BookPageRepository
}

//revive:disable:unexported-return // Intentionally returns an unexported struct to enforce dependency on the interface in other layers.
func NewBookPageUseCase(bookRepository port.BookRepository, bookPageRepository port.BookPageRepository) *bookPageUseCase {
	return &bookPageUseCase{
		BookRepository:     bookRepository,
		BookPageRepository: bookPageRepository,
	}
}

//revive:enable:unexported-return

func (u *bookPageUseCase) CreateBookPage(ctx context.Context, bookPage *entity.BookPage) (string, error) {
	_, err := u.BookRepository.GetBookByID(ctx, bookPage.BookID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch book: %w", err)
	}

	currentTime := time.Now()
	bookPage.CreatedAt = currentTime
	bookPage.UpdatedAt = currentTime

	newID, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate book page ID: %w", err)
	}
	bookPage.ID = newID.String()

	id, err := u.BookPageRepository.CreateBookPage(ctx, bookPage)
	if err != nil {
		return "", fmt.Errorf("failed to create book page: %w", err)
	}

	return id, nil
}

func (u *bookPageUseCase) GetAllBookPages(ctx context.Context, bookID string) ([]entity.BookPage, error) {
	bookPages, err := u.BookPageRepository.GetBookPagesByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages: %w", err)
	}

	return bookPages, nil
}

func (u *bookPageUseCase) GetBookPagesWithPagination(
	ctx context.Context,
	bookID string,
	pageSize int64,
	pageNumber int64,
) ([]entity.BookPage, error) {
	bookPages, err := u.BookPageRepository.GetBookpagesByBookIDWithPagination(ctx, bookID, pageSize, pageNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages: %w", err)
	}

	return bookPages, nil
}

func (u *bookPageUseCase) GetBookPagesByRange(
	ctx context.Context,
	bookID string,
	startPage int64,
	endPage int64,
) ([]entity.BookPage, error) {
	bookPages, err := u.BookPageRepository.GetBookpagesByPageRange(ctx, bookID, startPage, endPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages by range: %w", err)
	}

	return bookPages, nil
}

func (u *bookPageUseCase) GetBookPagesByOffset(
	ctx context.Context,
	bookID string,
	centerPage int64,
	offset int64,
) ([]entity.BookPage, error) {
	bookPages, err := u.BookPageRepository.GetBookpagesAroundPageNumber(ctx, bookID, centerPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages by offset: %w", err)
	}

	return bookPages, nil
}

func (u *bookPageUseCase) GetBookPageByID(ctx context.Context, id string) (entity.BookPage, error) {
	bookPage, err := u.BookPageRepository.GetBookPageByID(ctx, id)
	if err != nil {
		return entity.BookPage{}, fmt.Errorf("failed to fetch book page: %w", err)
	}

	return bookPage, nil
}
