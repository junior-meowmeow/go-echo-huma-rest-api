package handler

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

type ReviewHandler interface {
	PostReview(ctx context.Context, input *schema.ReviewInput) (*struct{}, error)
	GetReviews(ctx context.Context, _ *struct{}) (*schema.GetReviewsOutput, error)
}

type reviewHandler struct {
	ReviewUseCase usecase.ReviewUseCase
}

func NewReviewHandler(reviewUseCase usecase.ReviewUseCase) *reviewHandler {
	return &reviewHandler{
		ReviewUseCase: reviewUseCase,
	}
}

func (h *reviewHandler) PostReview(ctx context.Context, input *schema.ReviewInput) (*struct{}, error) {
	review := &entity.Review{
		Author:  input.Body.Author,
		Rating:  input.Body.Rating,
		Message: input.Body.Message,
	}

	if err := h.ReviewUseCase.PostReview(ctx, review); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *reviewHandler) GetReviews(ctx context.Context, _ *struct{}) (*schema.GetReviewsOutput, error) {
	reviews, err := h.ReviewUseCase.GetReviews(ctx, 100)
	if err != nil {
		return nil, err
	}

	reviewOutputs := convertReviews(reviews)

	resp := schema.GetReviewsOutput{
		Body: reviewOutputs,
	}

	return &resp, nil
}

func convertReviews(reviews []entity.Review) []schema.ReviewOutput {
	reviewOutputs := make([]schema.ReviewOutput, len(reviews))
	for i, r := range reviews {
		reviewOutputs[i] = convertReview(r)
	}
	return reviewOutputs
}

func convertReview(review entity.Review) schema.ReviewOutput {
	return schema.ReviewOutput{
		ID:        fmt.Sprint(review.ID),
		Author:    review.Author,
		Rating:    review.Rating,
		Message:   review.Message,
		CreatedAt: review.CreatedAt,
	}
}
