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

type bookRepository struct {
	db      *pgxpool.Pool
	queries *sqlcgen.Queries
}

//revive:disable:unexported-return
func NewBookRepository(db *pgxpool.Pool) *bookRepository {
	return &bookRepository{
		db:      db,
		queries: sqlcgen.New(db),
	}
}

//revive:enable:unexported-return

func (r *bookRepository) CreateBook(ctx context.Context, book *entity.Book) (string, error) {
	bookUUID, err := uuid.Parse(book.ID)
	if err != nil {
		return "", fmt.Errorf("invalid book ID format: %w", err)
	}

	err = r.queries.CreateBook(ctx, sqlcgen.CreateBookParams{
		ID:               bookUUID,
		Name:             book.Name,
		Description:      book.Description,
		Author:           book.Metadata.Author,
		Isbn:             book.Metadata.ISBN,
		Genre:            book.Metadata.Genre,
		CoverImageFileID: book.CoverImageFileID,
		CreatedAt:        book.CreatedAt,
		UpdatedAt:        book.UpdatedAt,
	})
	if err != nil {
		return "", fmt.Errorf("failed to insert book: %w", err)
	}

	return book.ID, nil
}

func (r *bookRepository) GetBookByID(ctx context.Context, id string) (entity.Book, error) {
	bookUUID, err := uuid.Parse(id)
	if err != nil {
		return entity.Book{}, fmt.Errorf("invalid book ID format: %w", err)
	}

	row, err := r.queries.GetBookByID(ctx, bookUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Book{}, fmt.Errorf("failed to get book: %w: %w", entity.ErrNotFound, err)
		}
		return entity.Book{}, err
	}

	return mapBookToEntity(row), nil
}

func (r *bookRepository) GetAllBooks(ctx context.Context) ([]entity.Book, error) {
	rows, err := r.queries.GetAllBooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books: %w", err)
	}

	books := make([]entity.Book, len(rows))
	for i, row := range rows {
		books[i] = mapBookToEntity(row)
	}

	return books, nil
}

func (r *bookRepository) GetBooksWithPagination(ctx context.Context, pageSize int64, pageNumber int64) ([]entity.Book, error) {
	skip := max((pageNumber-1)*pageSize, 0)

	rows, err := r.queries.GetBooksWithPagination(ctx, sqlcgen.GetBooksWithPaginationParams{
		Limit:  safeInt32(pageSize),
		Offset: safeInt32(skip),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books with pagination: %w", err)
	}

	books := make([]entity.Book, len(rows))
	for i, row := range rows {
		books[i] = mapBookToEntity(row)
	}

	return books, nil
}

// Helper function to map sqlc struct to domain entity.
func mapBookToEntity(row sqlcgen.Book) entity.Book {
	return entity.Book{
		ID:          row.ID.String(),
		Name:        row.Name,
		Description: row.Description,
		Metadata: entity.BookMetadata{
			Author: row.Author,
			ISBN:   row.Isbn,
			Genre:  row.Genre,
		},
		CoverImageFileID: row.CoverImageFileID,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
}
