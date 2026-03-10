package usecases

import (
	"context"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories/mongo_repositories"
)

type ReviewsUseCase interface {
	PostReview(ctx context.Context, review *entities.Review) error
	GetReviews(ctx context.Context, limit int64) ([]entities.Review, error)
}

type reviewsUseCase struct {
	ReviewsRepository mongo_repositories.ReviewsRepository
}

func NewReviewsUseCase(reviewsRepo mongo_repositories.ReviewsRepository) *reviewsUseCase {
	return &reviewsUseCase{
		ReviewsRepository: reviewsRepo,
	}
}

func (u *reviewsUseCase) PostReview(ctx context.Context, review *entities.Review) error {
	currentTime := time.Now()
	review.CreatedAt = currentTime

	if err := u.ReviewsRepository.CreateReview(ctx, review); err != nil {
		return err
	}

	return nil
}

func (u *reviewsUseCase) GetReviews(ctx context.Context, limit int64) ([]entities.Review, error) {
	records, err := u.ReviewsRepository.GetReviews(ctx, limit)
	if err != nil {
		return nil, err
	}

	return records, nil
}
