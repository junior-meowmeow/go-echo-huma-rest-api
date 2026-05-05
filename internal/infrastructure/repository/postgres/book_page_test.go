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
	"github.com/junior-meowmeow/go-echo-huma-rest-api/test/testhelper"
)

func TestPostgresBookPageRepository(t *testing.T) {
	dbPool := testhelper.SetupPostgresDatabase(t)
	ctx := context.Background()
	repo := postgres.NewBookPageRepository(dbPool)

	mockTime := time.Date(2025, 10, 25, 12, 0, 0, 0, time.UTC)

	// Helper function to insert a parent book to satisfy the foreign key constraint
	seedDummyBook := func(t *testing.T, bookID string) {
		t.Helper()
		query := `INSERT INTO books (id, name) VALUES ($1, $2)`
		_, err := dbPool.Exec(ctx, query, bookID, "Dummy Test Book")
		require.NoError(t, err, "failed to seed dummy book for foreign key constraint")
	}

	t.Run("CreateBookPage", func(t *testing.T) {
		testhelper.CleanPostgresTable(t, dbPool, "book_pages", "books")

		t.Run("Should create book page successfully", func(t *testing.T) {
			bookID := uuid.NewString()
			seedDummyBook(t, bookID)

			input := &entity.BookPage{
				ID:         uuid.NewString(),
				BookID:     bookID,
				PageNumber: 1,
				Content:    "Hello World",
				Metadata: entity.BookPageMetadata{
					IsBookmarked: true,
					Highlight:    "Important section",
				},
				AttachedImageFileID: "img-123",
				CreatedAt:           mockTime,
				UpdatedAt:           mockTime,
			}

			insertedID, err := repo.CreateBookPage(ctx, input)
			require.NoError(t, err)
			assert.Equal(t, input.ID, insertedID)

			// Verify insertion by querying the database directly
			var dbBookID, dbContent, dbHighlight, dbAttachedImage string
			var dbPageNumber int64
			var dbIsBookmarked bool

			query := `SELECT book_id, page_number, content, is_bookmarked, highlight, attached_image_file_id 
			          FROM book_pages WHERE id = $1`

			err = dbPool.QueryRow(ctx, query, insertedID).Scan(
				&dbBookID,
				&dbPageNumber,
				&dbContent,
				&dbIsBookmarked,
				&dbHighlight,
				&dbAttachedImage,
			)
			require.NoError(t, err, "failed to query raw database for created book page")

			assert.Equal(t, input.BookID, dbBookID)
			assert.Equal(t, input.PageNumber, dbPageNumber)
			assert.Equal(t, input.Content, dbContent)
			assert.Equal(t, input.Metadata.IsBookmarked, dbIsBookmarked)
			assert.Equal(t, input.Metadata.Highlight, dbHighlight)
			assert.Equal(t, input.AttachedImageFileID, dbAttachedImage)
		})
	})

	t.Run("GetBookPageByID", func(t *testing.T) {
		t.Run("Should return book page when exists", func(t *testing.T) {
			testhelper.CleanPostgresTable(t, dbPool, "book_pages", "books")

			bookID := uuid.NewString()
			seedDummyBook(t, bookID)

			expectedID := uuid.NewString()
			expectedContent := "Test Page Content"

			query := `INSERT INTO book_pages (id, book_id, page_number, content) VALUES ($1, $2, $3, $4)`
			_, err := dbPool.Exec(ctx, query, expectedID, bookID, 5, expectedContent)
			require.NoError(t, err, "failed to seed database with raw SQL")

			page, err := repo.GetBookPageByID(ctx, expectedID)

			require.NoError(t, err)
			assert.Equal(t, int64(5), page.PageNumber)
			assert.Equal(t, expectedContent, page.Content)
		})

		t.Run("Should return ErrNotFound when not exists", func(t *testing.T) {
			testhelper.CleanPostgresTable(t, dbPool, "book_pages", "books")

			_, err := repo.GetBookPageByID(ctx, uuid.NewString())
			require.ErrorIs(t, err, entity.ErrNotFound)
		})
	})

	t.Run("GetBookPagesByBookID", func(t *testing.T) {
		testhelper.CleanPostgresTable(t, dbPool, "book_pages", "books")

		bookID := uuid.NewString()
		seedDummyBook(t, bookID)

		seeds := []struct {
			ID         string
			PageNumber int64
		}{
			{ID: uuid.NewString(), PageNumber: 3},
			{ID: uuid.NewString(), PageNumber: 1},
			{ID: uuid.NewString(), PageNumber: 2},
		}

		query := `INSERT INTO book_pages (id, book_id, page_number) VALUES ($1, $2, $3)`
		for _, s := range seeds {
			_, err := dbPool.Exec(ctx, query, s.ID, bookID, s.PageNumber)
			require.NoError(t, err)
		}

		t.Run("Should return pages sorted by page number ascending", func(t *testing.T) {
			pages, err := repo.GetBookPagesByBookID(ctx, bookID)
			require.NoError(t, err)

			require.Len(t, pages, 3)
			assert.Equal(t, int64(1), pages[0].PageNumber)
			assert.Equal(t, int64(2), pages[1].PageNumber)
			assert.Equal(t, int64(3), pages[2].PageNumber)
		})
	})

	t.Run("GetBookpagesByBookIDWithPagination", func(t *testing.T) {
		testhelper.CleanPostgresTable(t, dbPool, "book_pages", "books")

		bookID := uuid.NewString()
		seedDummyBook(t, bookID)

		seeds := []struct {
			ID         string
			PageNumber int64
		}{
			{ID: uuid.NewString(), PageNumber: 1},
			{ID: uuid.NewString(), PageNumber: 2},
			{ID: uuid.NewString(), PageNumber: 3},
			{ID: uuid.NewString(), PageNumber: 4},
		}

		query := `INSERT INTO book_pages (id, book_id, page_number) VALUES ($1, $2, $3)`
		for _, s := range seeds {
			_, err := dbPool.Exec(ctx, query, s.ID, bookID, s.PageNumber)
			require.NoError(t, err)
		}

		t.Run("Should return book page 1-2", func(t *testing.T) {
			pages, err := repo.GetBookpagesByBookIDWithPagination(ctx, bookID, 2, 1)
			require.NoError(t, err)

			require.Len(t, pages, 2)
			assert.Equal(t, int64(1), pages[0].PageNumber)
			assert.Equal(t, int64(2), pages[1].PageNumber)
		})

		t.Run("Should return book page 3-4", func(t *testing.T) {
			pages, err := repo.GetBookpagesByBookIDWithPagination(ctx, bookID, 2, 2)
			require.NoError(t, err)

			require.Len(t, pages, 2)
			assert.Equal(t, int64(3), pages[0].PageNumber)
			assert.Equal(t, int64(4), pages[1].PageNumber)
		})
	})

	t.Run("GetBookpagesByPageRange", func(t *testing.T) {
		testhelper.CleanPostgresTable(t, dbPool, "book_pages", "books")

		bookID := uuid.NewString()
		seedDummyBook(t, bookID)

		seeds := []struct {
			ID         string
			PageNumber int64
		}{
			{ID: uuid.NewString(), PageNumber: 1},
			{ID: uuid.NewString(), PageNumber: 2},
			{ID: uuid.NewString(), PageNumber: 3},
			{ID: uuid.NewString(), PageNumber: 4},
		}

		query := `INSERT INTO book_pages (id, book_id, page_number) VALUES ($1, $2, $3)`
		for _, s := range seeds {
			_, err := dbPool.Exec(ctx, query, s.ID, bookID, s.PageNumber)
			require.NoError(t, err)
		}

		t.Run("Should return pages within range", func(t *testing.T) {
			pages, err := repo.GetBookpagesByPageRange(ctx, bookID, 2, 3)
			require.NoError(t, err)

			require.Len(t, pages, 2)
			assert.Equal(t, int64(2), pages[0].PageNumber)
			assert.Equal(t, int64(3), pages[1].PageNumber)
		})
	})

	t.Run("GetBookpagesAroundPageNumber", func(t *testing.T) {
		testhelper.CleanPostgresTable(t, dbPool, "book_pages", "books")

		bookID := uuid.NewString()
		seedDummyBook(t, bookID)

		seeds := []struct {
			ID         string
			PageNumber int64
		}{
			{ID: uuid.NewString(), PageNumber: 1},
			{ID: uuid.NewString(), PageNumber: 2},
			{ID: uuid.NewString(), PageNumber: 3},
			{ID: uuid.NewString(), PageNumber: 4},
			{ID: uuid.NewString(), PageNumber: 5},
		}

		query := `INSERT INTO book_pages (id, book_id, page_number) VALUES ($1, $2, $3)`
		for _, s := range seeds {
			_, err := dbPool.Exec(ctx, query, s.ID, bookID, s.PageNumber)
			require.NoError(t, err)
		}

		t.Run("Should return pages around center", func(t *testing.T) {
			pages, err := repo.GetBookpagesAroundPageNumber(ctx, bookID, 3, 1)
			require.NoError(t, err)

			require.Len(t, pages, 3)
			assert.Equal(t, int64(2), pages[0].PageNumber)
			assert.Equal(t, int64(3), pages[1].PageNumber)
			assert.Equal(t, int64(4), pages[2].PageNumber)
		})

		t.Run("Should handle correctly when offset is 0", func(t *testing.T) {
			pages, err := repo.GetBookpagesAroundPageNumber(ctx, bookID, 3, 0)
			require.NoError(t, err)

			require.Len(t, pages, 1)
			assert.Equal(t, int64(3), pages[0].PageNumber)
		})
	})
}
