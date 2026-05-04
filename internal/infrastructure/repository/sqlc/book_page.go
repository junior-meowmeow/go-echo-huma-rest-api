package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/sqlc/sqlcgen"
)

type bookPageRepository struct {
	db      *pgxpool.Pool
	queries *sqlcgen.Queries
}

//revive:disable:unexported-return
func NewBookPageRepository(db *pgxpool.Pool) *bookPageRepository {
	return &bookPageRepository{
		db:      db,
		queries: sqlcgen.New(db),
	}
}

//revive:enable:unexported-return

func (r *bookPageRepository) CreateBookPage(ctx context.Context, bookPage *entity.BookPage) (string, error) {
	bookPageUUID, err := uuid.Parse(bookPage.ID)
	if err != nil {
		return "", fmt.Errorf("invalid book page ID format: %w", err)
	}
	bookUUID, err := uuid.Parse(bookPage.BookID)
	if err != nil {
		return "", fmt.Errorf("invalid book ID format: %w", err)
	}

	err = r.queries.CreateBookPage(ctx, sqlcgen.CreateBookPageParams{
		ID:                  bookPageUUID,
		BookID:              bookUUID,
		PageNumber:          bookPage.PageNumber,
		Content:             bookPage.Content,
		IsBookmarked:        bookPage.Metadata.IsBookmarked,
		Highlight:           bookPage.Metadata.Highlight,
		AttachedImageFileID: bookPage.AttachedImageFileID,
		CreatedAt:           bookPage.CreatedAt,
		UpdatedAt:           bookPage.UpdatedAt,
	})
	if err != nil {
		return "", fmt.Errorf("failed to insert book page: %w", err)
	}

	return bookPage.ID, nil
}

func (r *bookPageRepository) GetBookPageByID(ctx context.Context, id string) (entity.BookPage, error) {
	bookPageUUID, err := uuid.Parse(id)
	if err != nil {
		return entity.BookPage{}, fmt.Errorf("invalid book page ID format: %w", err)
	}

	row, err := r.queries.GetBookPageByID(ctx, bookPageUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.BookPage{}, fmt.Errorf("failed to get book page: %w: %w", entity.ErrNotFound, err)
		}
		return entity.BookPage{}, err
	}

	return mapBookPageToEntity(row), nil
}

func (r *bookPageRepository) GetBookPagesByBookID(ctx context.Context, bookID string) ([]entity.BookPage, error) {
	bookUUID, err := uuid.Parse(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID format: %w", err)
	}

	rows, err := r.queries.GetBookPagesByBookID(ctx, bookUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages: %w", err)
	}

	pages := make([]entity.BookPage, len(rows))
	for i, row := range rows {
		pages[i] = mapBookPageToEntity(row)
	}

	return pages, nil
}

func (r *bookPageRepository) GetBookpagesByBookIDWithPagination(
	ctx context.Context,
	bookID string,
	pageSize int64,
	pageNumber int64,
) ([]entity.BookPage, error) {
	bookUUID, err := uuid.Parse(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID format: %w", err)
	}

	skip := max((pageNumber-1)*pageSize, 0)

	rows, err := r.queries.GetBookpagesByBookIDWithPagination(ctx, sqlcgen.GetBookpagesByBookIDWithPaginationParams{
		BookID: bookUUID,
		Limit:  safeInt32(pageSize),
		Offset: safeInt32(skip),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages with pagination: %w", err)
	}

	pages := make([]entity.BookPage, len(rows))
	for i, row := range rows {
		pages[i] = mapBookPageToEntity(row)
	}

	return pages, nil
}

func (r *bookPageRepository) GetBookpagesByPageRange(
	ctx context.Context,
	bookID string,
	startPage int64,
	endPage int64,
) ([]entity.BookPage, error) {
	bookUUID, err := uuid.Parse(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID format: %w", err)
	}

	rows, err := r.queries.GetBookpagesByPageRange(ctx, sqlcgen.GetBookpagesByPageRangeParams{
		BookID:       bookUUID,
		PageNumber:   startPage,
		PageNumber_2: endPage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages by range: %w", err)
	}

	pages := make([]entity.BookPage, len(rows))
	for i, row := range rows {
		pages[i] = mapBookPageToEntity(row)
	}

	return pages, nil
}

func (r *bookPageRepository) GetBookpagesAroundPageNumber(
	ctx context.Context,
	bookID string,
	centerPage int64,
	offset int64,
) ([]entity.BookPage, error) {
	bookUUID, err := uuid.Parse(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID format: %w", err)
	}

	startPage := max(centerPage-offset, 1)
	endPage := centerPage + offset

	rows, err := r.queries.GetBookpagesByPageRange(ctx, sqlcgen.GetBookpagesByPageRangeParams{
		BookID:       bookUUID,
		PageNumber:   startPage,
		PageNumber_2: endPage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book around page number: %w", err)
	}

	pages := make([]entity.BookPage, len(rows))
	for i, row := range rows {
		pages[i] = mapBookPageToEntity(row)
	}

	return pages, nil
}

// Helper function to map sqlc struct to domain entity.
func mapBookPageToEntity(row sqlcgen.BookPage) entity.BookPage {
	return entity.BookPage{
		ID:         row.ID.String(),
		BookID:     row.BookID.String(),
		PageNumber: row.PageNumber,
		Content:    row.Content,
		Metadata: entity.BookPageMetadata{
			IsBookmarked: row.IsBookmarked,
			Highlight:    row.Highlight,
		},
		AttachedImageFileID: row.AttachedImageFileID,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
	}
}
