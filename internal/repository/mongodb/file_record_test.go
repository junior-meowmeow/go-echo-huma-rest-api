package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository/mongodb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileRecordRepository(t *testing.T) {
	db := setupMongoDatabase(t)
	ctx := context.Background()

	repo := mongodb.NewFileRecordRepository(db)

	t.Run("Save and Get File Record", func(t *testing.T) {
		cleanCollection(t, repo.Collection)

		mockTime := time.Date(2025, 10, 25, 12, 0, 0, 0, time.UTC)
		inputRecord := &entity.FileRecord{
			FileName:    "test_image.png",
			Size:        1024,
			ContentType: "image/png",
			S3Key:       "uuid-test-key",
			CreatedAt:   mockTime,
		}

		// Test CreateFileRecord
		id, err := repo.CreateFileRecord(ctx, inputRecord)
		require.NoError(t, err, "CreateFileRecord should succeed")
		require.NotEmpty(t, id, "Returned ID should not be empty")

		// Test GetFileRecordByID
		fetchedRecord, err := repo.GetFileRecordByID(ctx, id)
		require.NoError(t, err, "GetFileRecordByID should succeed")

		// Assertions
		assert.Equal(t, inputRecord.FileName, fetchedRecord.FileName)
		assert.Equal(t, inputRecord.Size, fetchedRecord.Size)
		assert.Equal(t, inputRecord.S3Key, fetchedRecord.S3Key)
		assert.Equal(t, inputRecord.CreatedAt, fetchedRecord.CreatedAt)
	})

	t.Run("Get Non-Existent File", func(t *testing.T) {
		cleanCollection(t, repo.Collection)

		// Try to get a valid hex ID that doesn't exist
		fakeID := "000000000000000000000000"
		_, err := repo.GetFileRecordByID(ctx, fakeID)

		assert.Error(t, err)
		// Mongo driver returns "mongo: no documents in result"
		assert.Contains(t, err.Error(), "no documents")
	})

	t.Run("Get Invalid ID Format", func(t *testing.T) {
		cleanCollection(t, repo.Collection)

		_, err := repo.GetFileRecordByID(ctx, "invalid-hex-string")

		assert.Error(t, err)
		assert.Equal(t, "invalid ID format", err.Error())
	})
}
