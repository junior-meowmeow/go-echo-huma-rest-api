package database

import (
	"context"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongoAdapter) CleanFiles(t require.TestingT) {
	_, err := m.db.Collection("filerecords").DeleteMany(context.Background(), bson.D{})
	require.NoError(t, err)
}

func (m *MongoAdapter) CountFiles(t require.TestingT) int64 {
	count, err := m.db.Collection("filerecords").CountDocuments(context.Background(), bson.D{})
	require.NoError(t, err)
	return count
}
