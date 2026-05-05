package database

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongoAdapter) CleanUsers(t require.TestingT) {
	_, err := m.db.Collection("users").DeleteMany(context.Background(), bson.D{})
	require.NoError(t, err)
}

func (m *MongoAdapter) CountUsers(t require.TestingT) int64 {
	count, err := m.db.Collection("users").CountDocuments(context.Background(), bson.D{})
	require.NoError(t, err)
	return count
}

func (m *MongoAdapter) GetUserByUsername(t require.TestingT, username string) TestUserRecord {
	var doc bson.M

	err := m.db.Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(&doc)
	require.NoError(t, err)

	return TestUserRecord{
		ID:        fmt.Sprintf("%v", doc["_id"]),
		Username:  get[string](doc, "username"),
		Password:  get[string](doc, "password"),
		Role:      get[string](doc, "role"),
		CreatedAt: get[time.Time](doc, "createdAt"),
		UpdatedAt: get[time.Time](doc, "updatedAt"),
	}
}
