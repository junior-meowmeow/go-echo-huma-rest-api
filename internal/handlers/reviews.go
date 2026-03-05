package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/models"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories/mongo_repositories"
)

type ReviewsHandler interface {
	PostReview(ctx context.Context, input *models.ReviewInput) (*struct{}, error)
	GetReviews(ctx context.Context, _ *struct{}) (*models.GetReviewsOutput, error)
}

type reviewsHandler struct {
	ReviewsRepository mongo_repositories.ReviewsRepository
}

func NewReviewsHandler(reviewsRepo mongo_repositories.ReviewsRepository) *reviewsHandler {
	return &reviewsHandler{
		ReviewsRepository: reviewsRepo,
	}
}

func (h *reviewsHandler) PostReview(ctx context.Context, input *models.ReviewInput) (*struct{}, error) {
	record := &entities.Review{
		Author:    input.Body.Author,
		Rating:    input.Body.Rating,
		Message:   input.Body.Message,
		CreatedAt: time.Now(),
	}

	if err := h.ReviewsRepository.CreateReview(ctx, record); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *reviewsHandler) GetReviews(ctx context.Context, _ *struct{}) (*models.GetReviewsOutput, error) {
	records, err := h.ReviewsRepository.GetReviews(ctx, 100)
	if err != nil {
		return nil, err
	}

	results := make([]models.ReviewOutput, len(records))

	for i, r := range records {
		results[i] = models.ReviewOutput{
			ID:        fmt.Sprint(r.ID),
			Author:    r.Author,
			Rating:    r.Rating,
			Message:   r.Message,
			CreatedAt: r.CreatedAt,
		}
	}

	resp := &models.GetReviewsOutput{
		Body: results,
	}

	return resp, nil
}
