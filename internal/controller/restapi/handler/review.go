package handler

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

type ReviewHandler interface {
	CreateReview(ctx context.Context, request *schema.CreateReviewRequest) (*schema.CreateReviewResponse, error)
	GetReviews(ctx context.Context, request *schema.GetReviewsRequest) (*schema.GetReviewsResponse, error)
}

type reviewHandler struct {
	ReviewUseCase usecase.ReviewUseCase
}

//revive:disable:unexported-return // Intentionally returns an unexported struct to enforce dependency on the interface in other layers.
func NewReviewHandler(reviewUseCase usecase.ReviewUseCase) *reviewHandler {
	return &reviewHandler{
		ReviewUseCase: reviewUseCase,
	}
}

//revive:enable:unexported-return

func (h *reviewHandler) CreateReview(ctx context.Context, request *schema.CreateReviewRequest) (*schema.CreateReviewResponse, error) {
	review := &entity.Review{
		Author:  request.Body.Author,
		Rating:  request.Body.Rating,
		Message: request.Body.Message,
	}

	if err := h.ReviewUseCase.PostReview(ctx, review); err != nil {
		return nil, ResolveError(ctx, err)
	}

	return &schema.CreateReviewResponse{}, nil
}

func (h *reviewHandler) GetReviews(ctx context.Context, _ *schema.GetReviewsRequest) (*schema.GetReviewsResponse, error) {
	const reviewsLimit = 100
	reviews, err := h.ReviewUseCase.GetReviews(ctx, reviewsLimit)
	if err != nil {
		return nil, ResolveError(ctx, err)
	}

	reviewOutputs := mapEntityReviewsToSchema(reviews)

	resp := schema.GetReviewsResponse{}
	resp.Body.Data = reviewOutputs

	return &resp, nil
}

func mapEntityReviewsToSchema(reviews []entity.Review) []schema.Review {
	reviewOutputs := make([]schema.Review, len(reviews))
	for i, r := range reviews {
		reviewOutputs[i] = mapEntityReviewToSchema(r)
	}
	return reviewOutputs
}

func mapEntityReviewToSchema(review entity.Review) schema.Review {
	return schema.Review{
		ID:        review.ID,
		Author:    review.Author,
		Rating:    review.Rating,
		Message:   review.Message,
		CreatedAt: review.CreatedAt,
	}
}
