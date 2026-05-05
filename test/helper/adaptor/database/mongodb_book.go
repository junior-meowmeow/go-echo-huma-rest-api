package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongoAdapter) CleanBooks(t require.TestingT) {
	_, err := m.db.Collection("books").DeleteMany(context.Background(), bson.D{})
	require.NoError(t, err)
}

func (m *MongoAdapter) CountBooks(t require.TestingT) int64 {
	count, err := m.db.Collection("books").CountDocuments(context.Background(), bson.D{})
	require.NoError(t, err)
	return count
}

func (m *MongoAdapter) GetBookByID(t require.TestingT, id string) TestBookRecord {
	var doc bson.M
	bookUUID, err := uuid.Parse(id)
	require.NoError(t, err)

	err = m.db.Collection("books").FindOne(context.Background(), bson.M{"_id": bookUUID}).Decode(&doc)
	require.NoError(t, err)

	meta, _ := doc["metadata"].(bson.D)
	metaMap := make(map[string]any)
	for _, e := range meta {
		metaMap[e.Key] = e.Value
	}

	return TestBookRecord{
		ID:               id,
		Name:             get[string](doc, "name"),
		Description:      get[string](doc, "description"),
		Author:           get[string](metaMap, "author"),
		ISBN:             get[string](metaMap, "isbn"),
		Genre:            get[string](metaMap, "genre"),
		CoverImageFileID: getStringOptional(doc, "coverImageFileId"),
		CreatedAt:        get[time.Time](doc, "createdAt"),
		UpdatedAt:        get[time.Time](doc, "updatedAt"),
	}
}

func (m *MongoAdapter) SeedBooks(t require.TestingT, books []TestBookRecord) {
	docs := make([]any, len(books))
	for i, b := range books {
		bookUUID, _ := uuid.Parse(b.ID)
		docs[i] = bson.M{
			"_id":              bookUUID,
			"name":             b.Name,
			"description":      b.Description,
			"metadata":         bson.M{"author": b.Author, "isbn": b.ISBN, "genre": b.Genre},
			"coverImageFileId": b.CoverImageFileID,
			"createdAt":        b.CreatedAt,
			"updatedAt":        b.CreatedAt,
		}
	}
	_, err := m.db.Collection("books").InsertMany(context.Background(), docs)
	require.NoError(t, err)
}
