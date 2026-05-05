package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/postgres"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/test/helper/testenv"
)

func TestPostgresBookRepository(t *testing.T) {
	dbPool := testenv.SetupPostgresDatabase(t)
	ctx := context.Background()
	repo := postgres.NewBookRepository(dbPool)

	mockTime := time.Date(2025, 10, 25, 12, 0, 0, 0, time.UTC)

	t.Run("CreateBook", func(t *testing.T) {
		testenv.CleanPostgresTable(t, dbPool, "books")

		t.Run("Should create book successfully", func(t *testing.T) {
			input := &entity.Book{
				ID:          uuid.NewString(),
				Name:        "New Book",
				Description: "This is a book.",
				Metadata: entity.BookMetadata{
					Author: "The Author",
					ISBN:   "123-456",
					Genre:  "Fantasy",
				},
				CoverImageFileID: "file-123",
				CreatedAt:        mockTime,
				UpdatedAt:        mockTime,
			}

			insertedID, err := repo.CreateBook(ctx, input)
			require.NoError(t, err)
			assert.Equal(t, input.ID, insertedID)

			var dbName, dbDescription, dbAuthor, dbISBN, dbGenre, dbCoverImage string

			query := `SELECT name, description, author, isbn, genre, cover_image_file_id 
			          FROM books WHERE id = $1`

			err = dbPool.QueryRow(ctx, query, insertedID).Scan(
				&dbName,
				&dbDescription,
				&dbAuthor,
				&dbISBN,
				&dbGenre,
				&dbCoverImage,
			)
			require.NoError(t, err, "failed to query raw database for created book")

			assert.Equal(t, input.Name, dbName)
			assert.Equal(t, input.Description, dbDescription)
			assert.Equal(t, input.Metadata.Author, dbAuthor)
			assert.Equal(t, input.Metadata.ISBN, dbISBN)
			assert.Equal(t, input.Metadata.Genre, dbGenre)
			assert.Equal(t, input.CoverImageFileID, dbCoverImage)
		})
	})

	t.Run("GetBookByID", func(t *testing.T) {
		t.Run("Should return book when exists", func(t *testing.T) {
			testenv.CleanPostgresTable(t, dbPool, "books")

			expectedID := uuid.NewString()
			expectedName := "Test Book"

			query := `INSERT INTO books (id, name) VALUES ($1, $2)`
			_, err := dbPool.Exec(ctx, query, expectedID, expectedName)
			require.NoError(t, err, "failed to seed database with raw SQL")

			book, err := repo.GetBookByID(ctx, expectedID)

			require.NoError(t, err)
			assert.Equal(t, "Test Book", book.Name)
		})

		t.Run("Should return ErrNotFound when book does not exist", func(t *testing.T) {
			testenv.CleanPostgresTable(t, dbPool, "books")

			_, err := repo.GetBookByID(ctx, uuid.NewString())
			require.ErrorIs(t, err, entity.ErrNotFound)
		})
	})

	t.Run("GetAllBooks", func(t *testing.T) {
		testenv.CleanPostgresTable(t, dbPool, "books")

		seeds := []*entity.Book{
			{ID: uuid.NewString(), Name: "Book 1", CreatedAt: mockTime.Add(-2 * time.Hour)},
			{ID: uuid.NewString(), Name: "Book 2", CreatedAt: mockTime.Add(-1 * time.Hour)},
			{ID: uuid.NewString(), Name: "Book 3", CreatedAt: mockTime},
		}
		query := `INSERT INTO books (id, name, created_at) VALUES ($1, $2, $3)`
		for _, b := range seeds {
			_, err := dbPool.Exec(ctx, query, b.ID, b.Name, b.CreatedAt)
			require.NoError(t, err, "failed to seed database with raw SQL")
		}

		t.Run("Should return all books sorted by newest first", func(t *testing.T) {
			books, err := repo.GetAllBooks(ctx)
			require.NoError(t, err)

			require.Len(t, books, 3)
			assert.Equal(t, "Book 3", books[0].Name)
			assert.Equal(t, "Book 2", books[1].Name)
			assert.Equal(t, "Book 1", books[2].Name)
		})
	})

	t.Run("GetBooksWithPagination", func(t *testing.T) {
		testenv.CleanPostgresTable(t, dbPool, "books")

		seeds := []*entity.Book{
			{ID: uuid.NewString(), Name: "Book 1", CreatedAt: mockTime.Add(-4 * time.Hour)},
			{ID: uuid.NewString(), Name: "Book 2", CreatedAt: mockTime.Add(-3 * time.Hour)},
			{ID: uuid.NewString(), Name: "Book 3", CreatedAt: mockTime.Add(-2 * time.Hour)},
			{ID: uuid.NewString(), Name: "Book 4", CreatedAt: mockTime.Add(-1 * time.Hour)},
			{ID: uuid.NewString(), Name: "Book 5", CreatedAt: mockTime},
		}
		query := `INSERT INTO books (id, name, created_at) VALUES ($1, $2, $3)`
		for _, b := range seeds {
			_, err := dbPool.Exec(ctx, query, b.ID, b.Name, b.CreatedAt)
			require.NoError(t, err, "failed to seed database with raw SQL")
		}

		t.Run("Should return book 4-5", func(t *testing.T) {
			// Page 1, Size 2
			books, err := repo.GetBooksWithPagination(ctx, 2, 1)
			require.NoError(t, err)
			require.Len(t, books, 2)
			assert.Equal(t, "Book 5", books[0].Name)
			assert.Equal(t, "Book 4", books[1].Name)
		})

		t.Run("Should return book 2-3", func(t *testing.T) {
			// Page 2, Size 2
			books, err := repo.GetBooksWithPagination(ctx, 2, 2)
			require.NoError(t, err)
			require.Len(t, books, 2)
			assert.Equal(t, "Book 3", books[0].Name)
			assert.Equal(t, "Book 2", books[1].Name)
		})

		t.Run("Should return nothing when out of range", func(t *testing.T) {
			books, err := repo.GetBooksWithPagination(ctx, 2, 10)
			require.NoError(t, err)
			assert.Empty(t, books)
		})
	})
}
