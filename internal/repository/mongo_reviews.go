package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ReviewsRepository interface {
	CreateReview(ctx context.Context, review *ReviewRecord) error
	GetReviews(ctx context.Context, limit int64) ([]ReviewRecord, error)
}

type MongoReviewsRepository struct {
	Collection *mongo.Collection
}

func NewMongoReviewsRepository(db *mongo.Database) *MongoReviewsRepository {
	return &MongoReviewsRepository{
		Collection: db.Collection("reviews"),
	}
}

type ReviewRecord struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Author    string        `bson:"author"`
	Rating    int           `bson:"rating"`
	Message   string        `bson:"message"`
	CreatedAt time.Time     `bson:"createdAt"`
}

func (r *MongoReviewsRepository) CreateReview(ctx context.Context, record *ReviewRecord) error {
	_, err := r.Collection.InsertOne(ctx, record)
	return err
}

func (r *MongoReviewsRepository) GetReviews(ctx context.Context, limit int64) ([]ReviewRecord, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cur, err := r.Collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	results := make([]ReviewRecord, 0, limit)

	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
