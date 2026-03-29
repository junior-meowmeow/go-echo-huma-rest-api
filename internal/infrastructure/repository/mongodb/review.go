package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb/document"
)

type reviewRepository struct {
	Collection *mongo.Collection
}

//revive:disable:unexported-return // Intentionally returns an unexported struct to enforce dependency on the interface in other layers.
func NewReviewRepository(db *mongo.Database) *reviewRepository {
	return &reviewRepository{
		Collection: db.Collection("reviews"),
	}
}

//revive:enable:unexported-return

func (r *reviewRepository) CreateReview(ctx context.Context, review *entity.Review) error {
	doc, err := document.NewReviewDocument(review)
	if err != nil {
		return fmt.Errorf("failed to convert review to document: %w", err)
	}

	_, err = r.Collection.InsertOne(ctx, doc)
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

	var docs []document.ReviewDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode review documents: %w", err)
	}

	reviews := make([]entity.Review, len(docs))
	for i, doc := range docs {
		reviews[i] = doc.ToEntity()
	}

	return reviews, nil
}
