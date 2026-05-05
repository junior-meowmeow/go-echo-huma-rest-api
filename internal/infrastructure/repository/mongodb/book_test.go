package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb/document"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/test/testhelper"
)

func TestMongoBookRepository(t *testing.T) {
	db := testhelper.SetupMongoDatabase(t)
	ctx := context.Background()
	repo := mongodb.NewBookRepository(db)
	coll := repo.Collection

	mockTime := time.Date(2025, 10, 25, 12, 0, 0, 0, time.UTC)

	t.Run("CreateBook", func(t *testing.T) {
		testhelper.CleanMongoCollection(t, coll)

		t.Run("Should create book successfully", func(t *testing.T) {
			metadata := entity.BookMetadata{
				Author: "The Author",
			}

			input := &entity.Book{
				ID:          uuid.NewString(),
				Name:        "New Book",
				Description: "This is a book.",
				Metadata:    metadata,
			}

			insertedID, err := repo.CreateBook(ctx, input)
			require.NoError(t, err)
			assert.NotEmpty(t, insertedID)

			var doc document.BookDocument
			bookUUID, _ := document.StringToUUID(insertedID)
			err = coll.FindOne(ctx, bson.M{"_id": bookUUID}).Decode(&doc)
			require.NoError(t, err)
			assert.Equal(t, input.Name, doc.Name)
			assert.Equal(t, input.Description, doc.Description)
			assert.Equal(t, input.Metadata.Author, doc.Metadata.Author)
		})
	})

	t.Run("GetBookByID", func(t *testing.T) {
		t.Run("Should return book when exists", func(t *testing.T) {
			testhelper.CleanMongoCollection(t, coll)
			bookUUID, err := uuid.NewV7()
			require.NoError(t, err)
			seed := document.BookDocument{ID: bookUUID, Name: "Test Book"}
			_, err = coll.InsertOne(ctx, seed)
			require.NoError(t, err)

			book, err := repo.GetBookByID(ctx, bookUUID.String())

			require.NoError(t, err)
			assert.Equal(t, "Test Book", book.Name)
		})

		t.Run("Should return ErrNotFound when book does not exist", func(t *testing.T) {
			testhelper.CleanMongoCollection(t, coll)
			bookUUID, err := uuid.NewV7()
			require.NoError(t, err)

			_, err = repo.GetBookByID(ctx, bookUUID.String())
			require.ErrorIs(t, err, entity.ErrNotFound)
		})
	})

	t.Run("GetAllBooks", func(t *testing.T) {
		testhelper.CleanMongoCollection(t, coll)

		docs := []any{
			document.BookDocument{ID: uuid.New(), Name: "Book 1", CreatedAt: mockTime.Add(-2 * time.Hour)},
			document.BookDocument{ID: uuid.New(), Name: "Book 2", CreatedAt: mockTime.Add(-1 * time.Hour)},
			document.BookDocument{ID: uuid.New(), Name: "Book 3", CreatedAt: mockTime},
		}
		_, err := coll.InsertMany(ctx, docs)
		require.NoError(t, err)

		t.Run("Should return all books sorted by newest first", func(t *testing.T) {
			books, err := repo.GetAllBooks(ctx)
			require.NoError(t, err)
			require.Len(t, books, 3)
			assert.Equal(t, "Book 3", books[0].Name)
			assert.Equal(t, "Book 1", books[2].Name)
		})
	})

	t.Run("GetBooksWithPagination", func(t *testing.T) {
		testhelper.CleanMongoCollection(t, coll)

		docs := []any{
			document.BookDocument{ID: uuid.New(), Name: "Book 1", CreatedAt: mockTime.Add(-4 * time.Hour)},
			document.BookDocument{ID: uuid.New(), Name: "Book 2", CreatedAt: mockTime.Add(-3 * time.Hour)},
			document.BookDocument{ID: uuid.New(), Name: "Book 3", CreatedAt: mockTime.Add(-2 * time.Hour)},
			document.BookDocument{ID: uuid.New(), Name: "Book 4", CreatedAt: mockTime.Add(-1 * time.Hour)},
			document.BookDocument{ID: uuid.New(), Name: "Book 5", CreatedAt: mockTime},
		}
		_, err := coll.InsertMany(ctx, docs)
		require.NoError(t, err)

		t.Run("Should return book 1-2", func(t *testing.T) {
			books, err := repo.GetBooksWithPagination(ctx, 2, 1)
			require.NoError(t, err)
			require.Len(t, books, 2)
			assert.Equal(t, "Book 5", books[0].Name)
			assert.Equal(t, "Book 4", books[1].Name)
		})

		t.Run("Should return book 3-4", func(t *testing.T) {
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
