package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository/mongodb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReviewRepository(t *testing.T) {
	db := setupMongoDatabase(t)
	ctx := context.Background()

	repo := mongodb.NewReviewRepository(db)

	t.Run("Create and Get Single Review", func(t *testing.T) {
		cleanCollection(t, repo.Collection)

		mockTime := time.Date(2025, 10, 25, 12, 0, 0, 0, time.UTC)
		inputRecord := &entity.Review{
			Author:    "John Doe",
			Rating:    5,
			Message:   "Great service!",
			CreatedAt: mockTime,
		}

		// Test CreateReview
		err := repo.CreateReview(ctx, inputRecord)
		require.NoError(t, err, "CreateReview should succeed")

		// Test GetReviews
		reviews, err := repo.GetReviews(ctx, 10)
		require.NoError(t, err, "GetReviews should succeed")

		require.Len(t, reviews, 1)
		fetchedRecord := reviews[0]

		// Assertions
		assert.NotEmpty(t, fetchedRecord.ID, "ID should be generated")
		assert.Equal(t, inputRecord.Author, fetchedRecord.Author)
		assert.Equal(t, inputRecord.Rating, fetchedRecord.Rating)
		assert.Equal(t, inputRecord.Message, fetchedRecord.Message)
		assert.WithinDuration(t, inputRecord.CreatedAt, fetchedRecord.CreatedAt, time.Millisecond)
	})

	t.Run("Get Reviews Sorting and Limiting", func(t *testing.T) {
		cleanCollection(t, repo.Collection)

		mockTime1 := time.Date(2025, 10, 25, 11, 0, 0, 0, time.UTC)
		mockTime2 := time.Date(2025, 10, 25, 12, 0, 0, 0, time.UTC)
		mockTime3 := time.Date(2025, 10, 25, 13, 0, 0, 0, time.UTC)
		r1 := &entity.Review{
			Author:    "Oldest",
			Rating:    1,
			CreatedAt: mockTime1,
		}
		r2 := &entity.Review{
			Author:    "Middle",
			Rating:    3,
			CreatedAt: mockTime2,
		}
		r3 := &entity.Review{
			Author:    "Newest",
			Rating:    5,
			CreatedAt: mockTime3,
		}

		require.NoError(t, repo.CreateReview(ctx, r1))
		require.NoError(t, repo.CreateReview(ctx, r2))
		require.NoError(t, repo.CreateReview(ctx, r3))

		// Test Sorting (Should be Newest -> Middle -> Oldest)
		reviews, err := repo.GetReviews(ctx, 10)
		require.NoError(t, err)

		assert.True(t, len(reviews) == 3)

		// Verify descending order logic manually for safety
		for i := 0; i < len(reviews)-1; i++ {
			// Current record should be newer (after) or equal to the next record
			isNewerOrEqual := reviews[i].CreatedAt.After(reviews[i+1].CreatedAt) || reviews[i].CreatedAt.Equal(reviews[i+1].CreatedAt)
			assert.True(t, isNewerOrEqual, "Reviews should be sorted by CreatedAt desc")
		}

		// Test Limit
		limit := int64(2)
		limitedReviews, err := repo.GetReviews(ctx, limit)
		require.NoError(t, err)
		assert.Equal(t, int(limit), len(limitedReviews), "Should respect the limit")
	})
}
