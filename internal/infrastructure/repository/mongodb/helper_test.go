package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupMongoDatabase(t *testing.T) *mongo.Database {
	t.Helper()

	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:8.0")
	require.NoError(t, err, "failed to start container")

	t.Cleanup(func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	endpoint, err := mongodbContainer.ConnectionString(ctx)
	require.NoError(t, err, "failed to get connection string")

	client, err := mongo.Connect(options.Client().ApplyURI(endpoint))
	require.NoError(t, err, "failed to connect to MongoDB")

	err = client.Ping(ctx, nil)
	require.NoError(t, err, "failed to ping MongoDB")

	// Return a random DB for isolation
	return client.Database("test_db_" + t.Name())
}

func cleanCollection(t *testing.T, coll *mongo.Collection) {
	t.Helper()
	_, err := coll.DeleteMany(context.Background(), bson.D{})
	require.NoError(t, err, "failed to clear collection")
}
