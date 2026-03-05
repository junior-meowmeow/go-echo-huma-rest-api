package mongo_repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories/mongo_repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMongoFilesRepository(t *testing.T) {
	db := setupMongoDatabase(t)
	ctx := context.Background()

	repo := mongo_repositories.NewFileMetadataRepository(db)

	t.Run("Save and Get File Metadata", func(t *testing.T) {
		cleanCollection(t, repo.Collection)

		mockTime := time.Date(2025, 10, 25, 12, 0, 0, 0, time.UTC)
		inputRecord := &entities.FileMetadata{
			Filename:    "test_image.png",
			Size:        1024,
			ContentType: "image/png",
			S3Key:       "uuid-test-key",
			CreatedAt:   mockTime,
		}

		// Test SaveFileMetadata
		id, err := repo.SaveFileMetadata(ctx, inputRecord)
		require.NoError(t, err, "SaveFileMetadata should succeed")
		require.NotEmpty(t, id, "Returned ID should not be empty")

		// Test GetFileMetadataByID
		fetchedRecord, err := repo.GetFileMetadataByID(ctx, id)
		require.NoError(t, err, "GetFileMetadataByID should succeed")

		// Assertions
		assert.Equal(t, inputRecord.Filename, fetchedRecord.Filename)
		assert.Equal(t, inputRecord.Size, fetchedRecord.Size)
		assert.Equal(t, inputRecord.S3Key, fetchedRecord.S3Key)
		assert.Equal(t, inputRecord.CreatedAt, fetchedRecord.CreatedAt)
	})

	t.Run("Get Non-Existent File", func(t *testing.T) {
		cleanCollection(t, repo.Collection)

		// Try to get a valid hex ID that doesn't exist
		fakeID := "000000000000000000000000"
		_, err := repo.GetFileMetadataByID(ctx, fakeID)

		assert.Error(t, err)
		// Mongo driver returns "mongo: no documents in result"
		assert.Contains(t, err.Error(), "no documents")
	})

	t.Run("Get Invalid ID Format", func(t *testing.T) {
		cleanCollection(t, repo.Collection)

		_, err := repo.GetFileMetadataByID(ctx, "invalid-hex-string")

		assert.Error(t, err)
		assert.Equal(t, "invalid ID format", err.Error())
	})
}
