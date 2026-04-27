package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/port/mocks"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

func TestReviewUseCase(t *testing.T) {
	ctx := context.Background()

	t.Run("PostReview", func(t *testing.T) {
		t.Run("Should create review successfully", func(t *testing.T) {
			mockRepo := mocks.NewMockReviewRepository(t)

			input := &entity.Review{
				Message: "Great book!",
			}

			mockRepo.EXPECT().
				CreateReview(ctx, input).
				Run(func(_ context.Context, r *entity.Review) {
					// Verify timestamps
					assert.False(t, r.CreatedAt.IsZero())
					assert.False(t, r.UpdatedAt.IsZero())
					assert.WithinDuration(t, r.CreatedAt, r.UpdatedAt, time.Millisecond)
				}).
				Return(nil)

			uc := usecase.NewReviewUseCase(mockRepo)

			err := uc.PostReview(ctx, input)

			require.NoError(t, err)
		})

		t.Run("Should return error when repository fails", func(t *testing.T) {
			mockRepo := mocks.NewMockReviewRepository(t)

			input := &entity.Review{
				Message: "Bad book!",
			}

			mockErr := errors.New("db error")

			mockRepo.EXPECT().
				CreateReview(ctx, input).
				Return(mockErr)

			uc := usecase.NewReviewUseCase(mockRepo)

			err := uc.PostReview(ctx, input)

			require.Error(t, err)
			assert.ErrorIs(t, err, mockErr)
		})
	})

	t.Run("GetReviews", func(t *testing.T) {
		t.Run("Should return reviews successfully", func(t *testing.T) {
			mockRepo := mocks.NewMockReviewRepository(t)

			expected := []entity.Review{
				{Message: "Review 1"},
				{Message: "Review 2"},
			}

			mockRepo.EXPECT().
				GetReviews(ctx, int64(10)).
				Return(expected, nil)

			uc := usecase.NewReviewUseCase(mockRepo)

			result, err := uc.GetReviews(ctx, 10)

			require.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("Should return error when repository fails", func(t *testing.T) {
			mockRepo := mocks.NewMockReviewRepository(t)

			mockErr := errors.New("db error")

			mockRepo.EXPECT().
				GetReviews(ctx, int64(10)).
				Return(nil, mockErr)

			uc := usecase.NewReviewUseCase(mockRepo)

			result, err := uc.GetReviews(ctx, 10)

			require.Error(t, err)
			assert.Nil(t, result)
			assert.ErrorIs(t, err, mockErr)
		})
	})
}
