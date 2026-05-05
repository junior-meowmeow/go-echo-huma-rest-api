package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb/document"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/test/helper/testenv"
)

func TestMongoReviewRepository(t *testing.T) {
	db := testenv.SetupMongoDatabase(t)
	ctx := context.Background()
	repo := mongodb.NewReviewRepository(db)
	coll := repo.Collection

	mockTime := time.Date(2025, 10, 25, 12, 0, 0, 0, time.UTC)

	t.Run("CreateReview", func(t *testing.T) {
		testenv.CleanMongoCollection(t, coll)

		t.Run("Should create review successfully", func(t *testing.T) {
			input := &entity.Review{
				Author:    "Test User",
				Rating:    4,
				Message:   "Good Service!",
				CreatedAt: mockTime,
			}

			err := repo.CreateReview(ctx, input)
			require.NoError(t, err)

			var doc document.ReviewDocument
			err = coll.FindOne(ctx, bson.M{"author": "Test User"}).Decode(&doc)

			require.NoError(t, err, "Document should exist in MongoDB")
			assert.Equal(t, input.Rating, doc.Rating)
			assert.Equal(t, input.Message, doc.Message)
			assert.Equal(t, input.CreatedAt, doc.CreatedAt)
		})
	})

	t.Run("GetReviews", func(t *testing.T) {
		testenv.CleanMongoCollection(t, coll)

		docs := []any{
			document.ReviewDocument{Author: "Oldest", CreatedAt: mockTime.Add(-2 * time.Hour)},
			document.ReviewDocument{Author: "Middle", CreatedAt: mockTime.Add(-1 * time.Hour)},
			document.ReviewDocument{Author: "Newest", CreatedAt: mockTime},
		}
		_, err := coll.InsertMany(ctx, docs)
		require.NoError(t, err)

		t.Run("Should return reviews sorted by CreatedAt descending", func(t *testing.T) {
			reviews, err := repo.GetReviews(ctx, 10)

			require.NoError(t, err)
			require.Len(t, reviews, 3)
			assert.Equal(t, "Newest", reviews[0].Author)
			assert.Equal(t, "Middle", reviews[1].Author)
			assert.Equal(t, "Oldest", reviews[2].Author)
		})

		t.Run("Should limit number of reviews", func(t *testing.T) {
			limit := int64(2)
			reviews, err := repo.GetReviews(ctx, limit)

			require.NoError(t, err)
			assert.Len(t, reviews, int(limit))
			assert.Equal(t, "Newest", reviews[0].Author)
			assert.Equal(t, "Middle", reviews[1].Author)
		})
	})
}
