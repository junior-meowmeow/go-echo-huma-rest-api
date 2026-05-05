package database

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongoAdapter) CleanReviews(t require.TestingT) {
	_, err := m.db.Collection("reviews").DeleteMany(context.Background(), bson.D{})
	require.NoError(t, err)
}

func (m *MongoAdapter) CountReviews(t require.TestingT) int64 {
	count, err := m.db.Collection("reviews").CountDocuments(context.Background(), bson.D{})
	require.NoError(t, err)
	return count
}

func (m *MongoAdapter) GetAllReviews(t require.TestingT) []TestReviewRecord {
	cursor, err := m.db.Collection("reviews").Find(context.Background(), bson.D{})
	require.NoError(t, err)
	defer cursor.Close(context.Background())

	var records []TestReviewRecord
	for cursor.Next(context.Background()) {
		var doc bson.M
		err := cursor.Decode(&doc)
		require.NoError(t, err)

		records = append(records, TestReviewRecord{
			ID:        fmt.Sprintf("%v", doc["_id"]),
			Author:    get[string](doc, "author"),
			Rating:    getInt(doc, "rating"),
			Message:   get[string](doc, "message"),
			CreatedAt: get[time.Time](doc, "createdAt"),
			UpdatedAt: get[time.Time](doc, "updatedAt"),
		})
	}
	return records
}

func (m *MongoAdapter) SeedReviews(t require.TestingT, reviews []TestReviewRecord) {
	docs := make([]any, len(reviews))
	for i, r := range reviews {
		rating := r.Rating
		if rating > math.MaxInt32 {
			rating = math.MaxInt32
		} else if rating < math.MinInt32 {
			rating = math.MinInt32
		}
		docs[i] = bson.M{
			"author":    r.Author,
			"rating":    int32(rating),
			"message":   r.Message,
			"createdAt": r.CreatedAt,
			"updatedAt": r.UpdatedAt,
		}
	}
	_, err := m.db.Collection("reviews").InsertMany(context.Background(), docs)
	require.NoError(t, err)
}
