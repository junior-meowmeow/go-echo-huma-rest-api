package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongoAdapter) CleanBookPages(t require.TestingT) {
	_, err := m.db.Collection("book_pages").DeleteMany(context.Background(), bson.D{})
	require.NoError(t, err)
}

func (m *MongoAdapter) CountBookPages(t require.TestingT) int64 {
	count, err := m.db.Collection("book_pages").CountDocuments(context.Background(), bson.D{})
	require.NoError(t, err)
	return count
}

func (m *MongoAdapter) GetBookPageByID(t require.TestingT, id string) TestBookPageRecord {
	var doc bson.M
	pageUUID, err := uuid.Parse(id)
	require.NoError(t, err)

	err = m.db.Collection("book_pages").FindOne(context.Background(), bson.M{"_id": pageUUID}).Decode(&doc)
	require.NoError(t, err)

	meta, _ := doc["metadata"].(bson.D)
	metaMap := make(map[string]any)
	for _, e := range meta {
		metaMap[e.Key] = e.Value
	}

	bookIDBinary, _ := doc["book_id"].(bson.Binary)
	bookUUID, _ := uuid.FromBytes(bookIDBinary.Data)

	return TestBookPageRecord{
		ID:                  id,
		BookID:              bookUUID.String(),
		PageNumber:          get[int64](doc, "pageNumber"),
		Content:             get[string](doc, "content"),
		IsBookmarked:        get[bool](metaMap, "isBookmarked"),
		Highlight:           getStringOptional(metaMap, "highlight"),
		AttachedImageFileID: getStringOptional(doc, "attachedImageFileId"),
		CreatedAt:           get[time.Time](doc, "createdAt"),
		UpdatedAt:           get[time.Time](doc, "updatedAt"),
	}
}

func (m *MongoAdapter) SeedBookPages(t require.TestingT, pages []TestBookPageRecord) {
	docs := make([]any, len(pages))
	for i, p := range pages {
		pageUUID, _ := uuid.Parse(p.ID)
		bookUUID, _ := uuid.Parse(p.BookID)
		docs[i] = bson.M{
			"_id":                 pageUUID,
			"book_id":             bookUUID,
			"pageNumber":          p.PageNumber,
			"content":             p.Content,
			"metadata":            bson.M{"isBookmarked": p.IsBookmarked, "highlight": p.Highlight},
			"attachedImageFileId": p.AttachedImageFileID,
			"createdAt":           p.CreatedAt,
			"updatedAt":           p.UpdatedAt,
		}
	}
	_, err := m.db.Collection("book_pages").InsertMany(context.Background(), docs)
	require.NoError(t, err)
}
