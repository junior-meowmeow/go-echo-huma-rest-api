package usecase

import (
	"context"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository"
)

type ReviewUseCase interface {
	PostReview(ctx context.Context, review *entity.Review) error
	GetReviews(ctx context.Context, limit int64) ([]entity.Review, error)
}

type reviewUseCase struct {
	ReviewRepository repository.ReviewRepository
}

func NewReviewUseCase(reviewRepository repository.ReviewRepository) *reviewUseCase {
	return &reviewUseCase{
		ReviewRepository: reviewRepository,
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
	reviews, err := u.ReviewRepository.GetReviews(ctx, limit)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
