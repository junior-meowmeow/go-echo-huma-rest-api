package mongo_repositories

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories/mongo_repositories/documents"

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

func (r *reviewsRepository) CreateReview(ctx context.Context, review *entities.Review) error {
	document, err := documents.NewReviewDocument(review)
	if err != nil {
		return fmt.Errorf("failed to convert review to document: %w", err)
	}

	_, err = r.Collection.InsertOne(ctx, document)
	if err != nil {
		return fmt.Errorf("failed to insert review document: %w", err)
	}

	return nil
}

func (r *reviewsRepository) GetReviews(ctx context.Context, limit int64) ([]entities.Review, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.Collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documentsList []documents.ReviewDocument
	if err := cursor.All(ctx, &documentsList); err != nil {
		return nil, fmt.Errorf("failed to decode review documents: %w", err)
	}

	reviews := make([]entities.Review, len(documentsList))
	for i, document := range documentsList {
		reviews[i] = document.ToEntity()
	}

	return reviews, nil
}
