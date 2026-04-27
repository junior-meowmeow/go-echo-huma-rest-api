package handler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	usecasemocks "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase/mocks"
)

func TestReviewHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateReview", func(t *testing.T) {
		t.Run("Should create review successfully", func(t *testing.T) {
			mockUC := usecasemocks.NewMockReviewUseCase(t)

			req := &schema.CreateReviewRequest{}
			req.Body.Author = "John"
			req.Body.Rating = 5
			req.Body.Message = "Great!"

			mockUC.EXPECT().
				PostReview(mock.Anything, mock.MatchedBy(func(r *entity.Review) bool {
					return r.Author == "John" &&
						r.Rating == 5 &&
						r.Message == "Great!"
				})).
				Return(nil)

			h := handler.NewReviewHandler(mockUC)

			resp, err := h.CreateReview(ctx, req)

			require.NoError(t, err)
			assert.NotNil(t, resp)
		})

		t.Run("Should return error when usecase fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockReviewUseCase(t)

			mockErr := errors.New("db error")

			req := &schema.CreateReviewRequest{}

			mockUC.EXPECT().
				PostReview(mock.Anything, mock.Anything).
				Return(mockErr)

			h := handler.NewReviewHandler(mockUC)

			resp, err := h.CreateReview(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)

			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})

	t.Run("GetReviews", func(t *testing.T) {
		t.Run("Should return mapped reviews", func(t *testing.T) {
			mockUC := usecasemocks.NewMockReviewUseCase(t)

			now := time.Now()

			input := []entity.Review{
				{
					ID:        "1",
					Author:    "Alice",
					Rating:    5,
					Message:   "Excellent",
					CreatedAt: now,
				},
				{
					ID:        "2",
					Author:    "Bob",
					Rating:    4,
					Message:   "Good",
					CreatedAt: now,
				},
			}

			mockUC.EXPECT().
				GetReviews(mock.Anything, int64(100)).
				Return(input, nil)

			h := handler.NewReviewHandler(mockUC)

			resp, err := h.GetReviews(ctx, &schema.GetReviewsRequest{})

			require.NoError(t, err)
			require.NotNil(t, resp)

			require.Len(t, resp.Body.Data, 2)

			assert.Equal(t, "1", resp.Body.Data[0].ID)
			assert.Equal(t, "Alice", resp.Body.Data[0].Author)
			assert.Equal(t, 5, resp.Body.Data[0].Rating)
			assert.Equal(t, "Excellent", resp.Body.Data[0].Message)
			assert.WithinDuration(t, now, resp.Body.Data[0].CreatedAt, time.Millisecond)

			assert.Equal(t, "2", resp.Body.Data[1].ID)
			assert.Equal(t, "Bob", resp.Body.Data[1].Author)
		})

		t.Run("Should return empty list when no reviews", func(t *testing.T) {
			mockUC := usecasemocks.NewMockReviewUseCase(t)

			mockUC.EXPECT().
				GetReviews(mock.Anything, int64(100)).
				Return([]entity.Review{}, nil)

			h := handler.NewReviewHandler(mockUC)

			resp, err := h.GetReviews(ctx, &schema.GetReviewsRequest{})

			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Empty(t, resp.Body.Data)
		})

		t.Run("Should return error when usecase fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockReviewUseCase(t)

			mockErr := errors.New("db error")

			mockUC.EXPECT().
				GetReviews(mock.Anything, int64(100)).
				Return(nil, mockErr)

			h := handler.NewReviewHandler(mockUC)

			resp, err := h.GetReviews(ctx, &schema.GetReviewsRequest{})

			require.Error(t, err)
			assert.Nil(t, resp)

			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})
}
