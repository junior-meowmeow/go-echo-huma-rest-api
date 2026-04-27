package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb/document"
)

func TestBookPageRepository(t *testing.T) {
	db := setupMongoDatabase(t)
	ctx := context.Background()
	repo := mongodb.NewBookPageRepository(db)
	coll := repo.Collection

	t.Run("CreateBookPage", func(t *testing.T) {
		cleanCollection(t, coll)

		t.Run("Should create book page successfully", func(t *testing.T) {
			bookID := bson.NewObjectID()

			input := &entity.BookPage{
				BookID:     bookID.Hex(),
				PageNumber: 1,
				Content:    "Hello World",
			}

			insertedID, err := repo.CreateBookPage(ctx, input)
			require.NoError(t, err)
			assert.NotEmpty(t, insertedID)

			var doc document.BookPageDocument
			oid, _ := document.StringToObjectID(insertedID)

			err = coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
			require.NoError(t, err)

			assert.Equal(t, input.PageNumber, doc.PageNumber)
			assert.Equal(t, input.Content, doc.Content)
		})
	})

	t.Run("GetBookPageByID", func(t *testing.T) {
		t.Run("Should return book page when exists", func(t *testing.T) {
			cleanCollection(t, coll)

			oid := bson.NewObjectID()
			seed := document.BookPageDocument{
				ID:         oid,
				PageNumber: 5,
				Content:    "Test Page",
			}

			_, err := coll.InsertOne(ctx, seed)
			require.NoError(t, err)

			idStr, _ := document.IDToString(oid)
			page, err := repo.GetBookPageByID(ctx, idStr)

			require.NoError(t, err)
			assert.Equal(t, int64(5), page.PageNumber)
			assert.Equal(t, "Test Page", page.Content)
		})

		t.Run("Should return ErrNotFound when not exists", func(t *testing.T) {
			cleanCollection(t, coll)

			oid := bson.NewObjectID()
			idStr, _ := document.IDToString(oid)

			_, err := repo.GetBookPageByID(ctx, idStr)
			require.ErrorIs(t, err, entity.ErrNotFound)
		})
	})

	t.Run("GetBookPagesByBookID", func(t *testing.T) {
		cleanCollection(t, coll)

		bookOID := bson.NewObjectID()

		docs := []any{
			document.BookPageDocument{BookID: bookOID, PageNumber: 3},
			document.BookPageDocument{BookID: bookOID, PageNumber: 1},
			document.BookPageDocument{BookID: bookOID, PageNumber: 2},
		}

		_, err := coll.InsertMany(ctx, docs)
		require.NoError(t, err)

		t.Run("Should return pages sorted by page number ascending", func(t *testing.T) {
			pages, err := repo.GetBookPagesByBookID(ctx, bookOID.Hex())
			require.NoError(t, err)

			require.Len(t, pages, 3)
			assert.Equal(t, int64(1), pages[0].PageNumber)
			assert.Equal(t, int64(3), pages[2].PageNumber)
		})
	})

	t.Run("GetBookpagesByBookIDWithPagination", func(t *testing.T) {
		cleanCollection(t, coll)

		bookOID := bson.NewObjectID()

		docs := []any{
			document.BookPageDocument{BookID: bookOID, PageNumber: 1},
			document.BookPageDocument{BookID: bookOID, PageNumber: 2},
			document.BookPageDocument{BookID: bookOID, PageNumber: 3},
			document.BookPageDocument{BookID: bookOID, PageNumber: 4},
		}

		_, err := coll.InsertMany(ctx, docs)
		require.NoError(t, err)

		t.Run("Should return book page 1-2", func(t *testing.T) {
			pages, err := repo.GetBookpagesByBookIDWithPagination(ctx, bookOID.Hex(), 2, 1)
			require.NoError(t, err)

			require.Len(t, pages, 2)
			assert.Equal(t, int64(1), pages[0].PageNumber)
			assert.Equal(t, int64(2), pages[1].PageNumber)
		})

		t.Run("Should return book page 3-4", func(t *testing.T) {
			pages, err := repo.GetBookpagesByBookIDWithPagination(ctx, bookOID.Hex(), 2, 2)
			require.NoError(t, err)

			require.Len(t, pages, 2)
			assert.Equal(t, int64(3), pages[0].PageNumber)
			assert.Equal(t, int64(4), pages[1].PageNumber)
		})
	})

	t.Run("GetBookpagesByPageRange", func(t *testing.T) {
		cleanCollection(t, coll)

		bookOID := bson.NewObjectID()

		docs := []any{
			document.BookPageDocument{BookID: bookOID, PageNumber: 1},
			document.BookPageDocument{BookID: bookOID, PageNumber: 2},
			document.BookPageDocument{BookID: bookOID, PageNumber: 3},
			document.BookPageDocument{BookID: bookOID, PageNumber: 4},
		}

		_, err := coll.InsertMany(ctx, docs)
		require.NoError(t, err)

		t.Run("Should return pages within range", func(t *testing.T) {
			pages, err := repo.GetBookpagesByPageRange(ctx, bookOID.Hex(), 2, 3)
			require.NoError(t, err)

			require.Len(t, pages, 2)
			assert.Equal(t, int64(2), pages[0].PageNumber)
			assert.Equal(t, int64(3), pages[1].PageNumber)
		})
	})

	t.Run("GetBookpagesAroundPageNumber", func(t *testing.T) {
		cleanCollection(t, coll)

		bookOID := bson.NewObjectID()

		docs := []any{
			document.BookPageDocument{BookID: bookOID, PageNumber: 1},
			document.BookPageDocument{BookID: bookOID, PageNumber: 2},
			document.BookPageDocument{BookID: bookOID, PageNumber: 3},
			document.BookPageDocument{BookID: bookOID, PageNumber: 4},
			document.BookPageDocument{BookID: bookOID, PageNumber: 5},
		}

		_, err := coll.InsertMany(ctx, docs)
		require.NoError(t, err)

		t.Run("Should return pages around center", func(t *testing.T) {
			pages, err := repo.GetBookpagesAroundPageNumber(ctx, bookOID.Hex(), 3, 1)
			require.NoError(t, err)

			require.Len(t, pages, 3)
			assert.Equal(t, int64(2), pages[0].PageNumber)
			assert.Equal(t, int64(3), pages[1].PageNumber)
			assert.Equal(t, int64(4), pages[2].PageNumber)
		})

		t.Run("Should handle correctly when offset is 0", func(t *testing.T) {
			pages, err := repo.GetBookpagesAroundPageNumber(ctx, bookOID.Hex(), 3, 0)
			require.NoError(t, err)

			require.Len(t, pages, 1)
			assert.Equal(t, int64(3), pages[0].PageNumber)
		})
	})
}
