package handler

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

type ReviewHandler interface {
	CreateReview(ctx context.Context, request *schema.CreateReviewRequest) (*schema.CreateReviewResponse, error)
	GetReviews(ctx context.Context, request *schema.GetReviewsRequest) (*schema.GetReviewsResponse, error)
}

type reviewHandler struct {
	ReviewUseCase usecase.ReviewUseCase
}

func NewReviewHandler(reviewUseCase usecase.ReviewUseCase) *reviewHandler {
	return &reviewHandler{
		ReviewUseCase: reviewUseCase,
	}
}

func (h *reviewHandler) CreateReview(ctx context.Context, request *schema.CreateReviewRequest) (*schema.CreateReviewResponse, error) {
	review := &entity.Review{
		Author:  request.Body.Author,
		Rating:  request.Body.Rating,
		Message: request.Body.Message,
	}

	if err := h.ReviewUseCase.PostReview(ctx, review); err != nil {
		return nil, err
	}

	return &schema.CreateReviewResponse{}, nil
}

func (h *reviewHandler) GetReviews(ctx context.Context, _ *schema.GetReviewsRequest) (*schema.GetReviewsResponse, error) {
	reviews, err := h.ReviewUseCase.GetReviews(ctx, 100)
	if err != nil {
		return nil, err
	}

	reviewOutputs := convertReviews(reviews)

	resp := schema.GetReviewsResponse{}
	resp.Body.Data = reviewOutputs

	return &resp, nil
}

func convertReviews(reviews []entity.Review) []schema.Review {
	reviewOutputs := make([]schema.Review, len(reviews))
	for i, r := range reviews {
		reviewOutputs[i] = convertReview(r)
	}
	return reviewOutputs
}

func convertReview(review entity.Review) schema.Review {
	return schema.Review{
		ID:        review.ID,
		Author:    review.Author,
		Rating:    review.Rating,
		Message:   review.Message,
		CreatedAt: review.CreatedAt,
	}
}
