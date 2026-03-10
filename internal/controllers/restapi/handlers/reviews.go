package handlers

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controllers/restapi/models"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecases"
)

type ReviewsHandler interface {
	PostReview(ctx context.Context, input *models.ReviewInput) (*struct{}, error)
	GetReviews(ctx context.Context, _ *struct{}) (*models.GetReviewsOutput, error)
}

type reviewsHandler struct {
	ReviewsUseCase usecases.ReviewsUseCase
}

func NewReviewsHandler(reviewsUseCase usecases.ReviewsUseCase) *reviewsHandler {
	return &reviewsHandler{
		ReviewsUseCase: reviewsUseCase,
	}
}

func (h *reviewsHandler) PostReview(ctx context.Context, input *models.ReviewInput) (*struct{}, error) {
	review := &entities.Review{
		Author:  input.Body.Author,
		Rating:  input.Body.Rating,
		Message: input.Body.Message,
	}

	if err := h.ReviewsUseCase.PostReview(ctx, review); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *reviewsHandler) GetReviews(ctx context.Context, _ *struct{}) (*models.GetReviewsOutput, error) {
	reviews, err := h.ReviewsUseCase.GetReviews(ctx, 100)
	if err != nil {
		return nil, err
	}

	reviewOutputs := convertReviews(reviews)

	resp := models.GetReviewsOutput{
		Body: reviewOutputs,
	}

	return &resp, nil
}

func convertReviews(reviews []entities.Review) []models.ReviewOutput {
	reviewOutputs := make([]models.ReviewOutput, len(reviews))
	for i, r := range reviews {
		reviewOutputs[i] = convertReview(r)
	}
	return reviewOutputs
}

func convertReview(review entities.Review) models.ReviewOutput {
	return models.ReviewOutput{
		ID:        fmt.Sprint(review.ID),
		Author:    review.Author,
		Rating:    review.Rating,
		Message:   review.Message,
		CreatedAt: review.CreatedAt,
	}
}
