package usecase

import (
	"context"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/port"
)

type ReviewUseCase interface {
	PostReview(ctx context.Context, review *entity.Review) error
	GetReviews(ctx context.Context, limit int64) ([]entity.Review, error)
}

type reviewUseCase struct {
	ReviewRepository port.ReviewRepository
}

//revive:disable:unexported-return // Intentionally returns an unexported struct to enforce dependency on the interface in other layers.
func NewReviewUseCase(reviewRepository port.ReviewRepository) *reviewUseCase {
	return &reviewUseCase{
		ReviewRepository: reviewRepository,
	}
}

//revive:enable:unexported-return

func (u *reviewUseCase) PostReview(ctx context.Context, review *entity.Review) error {
	currentTime := time.Now()
	review.CreatedAt = currentTime
	review.UpdatedAt = currentTime

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
