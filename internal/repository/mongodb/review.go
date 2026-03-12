package mongodb

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository/mongodb/document"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ReviewRepository interface {
	CreateReview(ctx context.Context, review *entity.Review) error
	GetReviews(ctx context.Context, limit int64) ([]entity.Review, error)
}

type reviewRepository struct {
	Collection *mongo.Collection
}

func NewReviewRepository(db *mongo.Database) *reviewRepository {
	return &reviewRepository{
		Collection: db.Collection("reviews"),
	}
}

func (r *reviewRepository) CreateReview(ctx context.Context, review *entity.Review) error {
	document, err := document.NewReviewDocument(review)
	if err != nil {
		return fmt.Errorf("failed to convert review to document: %w", err)
	}

	_, err = r.Collection.InsertOne(ctx, document)
	if err != nil {
		return fmt.Errorf("failed to insert review document: %w", err)
	}

	return nil
}

func (r *reviewRepository) GetReviews(ctx context.Context, limit int64) ([]entity.Review, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.Collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documentsList []document.ReviewDocument
	if err := cursor.All(ctx, &documentsList); err != nil {
		return nil, fmt.Errorf("failed to decode review documents: %w", err)
	}

	reviews := make([]entity.Review, len(documentsList))
	for i, document := range documentsList {
		reviews[i] = document.ToEntity()
	}

	return reviews, nil
}
