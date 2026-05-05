package mongodb_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb/document"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/test/testhelper"
)

func TestMongoBookPageRepository(t *testing.T) {
	db := testhelper.SetupMongoDatabase(t)
	ctx := context.Background()
	repo := mongodb.NewBookPageRepository(db)
	coll := repo.Collection

	t.Run("CreateBookPage", func(t *testing.T) {
		testhelper.CleanMongoCollection(t, coll)

		t.Run("Should create book page successfully", func(t *testing.T) {
			input := &entity.BookPage{
				ID:         uuid.NewString(),
				PageNumber: 1,
				Content:    "Hello World",
			}

			insertedID, err := repo.CreateBookPage(ctx, input)
			require.NoError(t, err)
			assert.NotEmpty(t, insertedID)

			var doc document.BookPageDocument
			bookPageUUID, _ := document.StringToUUID(insertedID)

			err = coll.FindOne(ctx, bson.M{"_id": bookPageUUID}).Decode(&doc)
			require.NoError(t, err)

			assert.Equal(t, input.PageNumber, doc.PageNumber)
			assert.Equal(t, input.Content, doc.Content)
		})
	})

	t.Run("GetBookPageByID", func(t *testing.T) {
		t.Run("Should return book page when exists", func(t *testing.T) {
			testhelper.CleanMongoCollection(t, coll)

			bookPageUUID, err := uuid.NewV7()
			require.NoError(t, err)
			seed := document.BookPageDocument{
				ID:         bookPageUUID,
				PageNumber: 5,
				Content:    "Test Page",
			}

			_, err = coll.InsertOne(ctx, seed)
			require.NoError(t, err)

			page, err := repo.GetBookPageByID(ctx, bookPageUUID.String())

			require.NoError(t, err)
			assert.Equal(t, int64(5), page.PageNumber)
			assert.Equal(t, "Test Page", page.Content)
		})

		t.Run("Should return ErrNotFound when not exists", func(t *testing.T) {
			testhelper.CleanMongoCollection(t, coll)

			bookPageUUID, err := uuid.NewV7()
			require.NoError(t, err)

			_, err = repo.GetBookPageByID(ctx, bookPageUUID.String())
			require.ErrorIs(t, err, entity.ErrNotFound)
		})
	})

	t.Run("GetBookPagesByBookID", func(t *testing.T) {
		testhelper.CleanMongoCollection(t, coll)

		bookID, err := uuid.NewV7()
		require.NoError(t, err)

		docs := []any{
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 3},
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 1},
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 2},
		}

		_, err = coll.InsertMany(ctx, docs)
		require.NoError(t, err)

		t.Run("Should return pages sorted by page number ascending", func(t *testing.T) {
			pages, err := repo.GetBookPagesByBookID(ctx, bookID.String())
			require.NoError(t, err)

			require.Len(t, pages, 3)
			assert.Equal(t, int64(1), pages[0].PageNumber)
			assert.Equal(t, int64(3), pages[2].PageNumber)
		})
	})

	t.Run("GetBookpagesByBookIDWithPagination", func(t *testing.T) {
		testhelper.CleanMongoCollection(t, coll)

		bookID, err := uuid.NewV7()
		require.NoError(t, err)

		docs := []any{
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 1},
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 2},
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 3},
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 4},
		}

		_, err = coll.InsertMany(ctx, docs)
		require.NoError(t, err)

		t.Run("Should return book page 1-2", func(t *testing.T) {
			pages, err := repo.GetBookpagesByBookIDWithPagination(ctx, bookID.String(), 2, 1)
			require.NoError(t, err)

			require.Len(t, pages, 2)
			assert.Equal(t, int64(1), pages[0].PageNumber)
			assert.Equal(t, int64(2), pages[1].PageNumber)
		})

		t.Run("Should return book page 3-4", func(t *testing.T) {
			pages, err := repo.GetBookpagesByBookIDWithPagination(ctx, bookID.String(), 2, 2)
			require.NoError(t, err)

			require.Len(t, pages, 2)
			assert.Equal(t, int64(3), pages[0].PageNumber)
			assert.Equal(t, int64(4), pages[1].PageNumber)
		})
	})

	t.Run("GetBookpagesByPageRange", func(t *testing.T) {
		testhelper.CleanMongoCollection(t, coll)

		bookID, err := uuid.NewV7()
		require.NoError(t, err)

		docs := []any{
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 1},
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 2},
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 3},
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 4},
		}

		_, err = coll.InsertMany(ctx, docs)
		require.NoError(t, err)

		t.Run("Should return pages within range", func(t *testing.T) {
			pages, err := repo.GetBookpagesByPageRange(ctx, bookID.String(), 2, 3)
			require.NoError(t, err)

			require.Len(t, pages, 2)
			assert.Equal(t, int64(2), pages[0].PageNumber)
			assert.Equal(t, int64(3), pages[1].PageNumber)
		})
	})

	t.Run("GetBookpagesAroundPageNumber", func(t *testing.T) {
		testhelper.CleanMongoCollection(t, coll)

		bookID, err := uuid.NewV7()
		require.NoError(t, err)

		docs := []any{
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 1},
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 2},
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 3},
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 4},
			document.BookPageDocument{ID: uuid.New(), BookID: bookID, PageNumber: 5},
		}

		_, err = coll.InsertMany(ctx, docs)
		require.NoError(t, err)

		t.Run("Should return pages around center", func(t *testing.T) {
			pages, err := repo.GetBookpagesAroundPageNumber(ctx, bookID.String(), 3, 1)
			require.NoError(t, err)

			require.Len(t, pages, 3)
			assert.Equal(t, int64(2), pages[0].PageNumber)
			assert.Equal(t, int64(3), pages[1].PageNumber)
			assert.Equal(t, int64(4), pages[2].PageNumber)
		})

		t.Run("Should handle correctly when offset is 0", func(t *testing.T) {
			pages, err := repo.GetBookpagesAroundPageNumber(ctx, bookID.String(), 3, 0)
			require.NoError(t, err)

			require.Len(t, pages, 1)
			assert.Equal(t, int64(3), pages[0].PageNumber)
		})
	})
}
