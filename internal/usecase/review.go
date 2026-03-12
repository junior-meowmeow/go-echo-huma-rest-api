package usecase

import (
	"context"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository/mongodb"
)

type ReviewUseCase interface {
	PostReview(ctx context.Context, review *entity.Review) error
	GetReviews(ctx context.Context, limit int64) ([]entity.Review, error)
}

type reviewUseCase struct {
	ReviewRepository mongodb.ReviewRepository
}

func NewReviewUseCase(reviewsRepo mongodb.ReviewRepository) *reviewUseCase {
	return &reviewUseCase{
		ReviewRepository: reviewsRepo,
	}
}

func (u *reviewUseCase) PostReview(ctx context.Context, review *entity.Review) error {
	currentTime := time.Now()
	review.CreatedAt = currentTime

	if err := u.ReviewRepository.CreateReview(ctx, review); err != nil {
		return err
	}

	return nil
}

func (u *reviewUseCase) GetReviews(ctx context.Context, limit int64) ([]entity.Review, error) {
	records, err := u.ReviewRepository.GetReviews(ctx, limit)
	if err != nil {
		return nil, err
	}

	return records, nil
}
