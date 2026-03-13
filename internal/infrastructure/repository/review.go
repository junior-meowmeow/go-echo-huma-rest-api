package repository

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
)

type ReviewRepository interface {
	CreateReview(ctx context.Context, review *entity.Review) error
	GetReviews(ctx context.Context, limit int64) ([]entity.Review, error)
}
