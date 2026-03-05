package mongo_repositories

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ReviewsRepository interface {
	CreateReview(ctx context.Context, review *entities.Review) error
	GetReviews(ctx context.Context, limit int64) ([]entities.Review, error)
}

type reviewsRepository struct {
	Collection *mongo.Collection
}

func NewReviewsRepository(db *mongo.Database) *reviewsRepository {
	return &reviewsRepository{
		Collection: db.Collection("reviews"),
	}
}

func (r *reviewsRepository) CreateReview(ctx context.Context, record *entities.Review) error {
	_, err := r.Collection.InsertOne(ctx, record)
	return err
}

func (r *reviewsRepository) GetReviews(ctx context.Context, limit int64) ([]entities.Review, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cur, err := r.Collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	results := make([]entities.Review, 0, limit)

	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
