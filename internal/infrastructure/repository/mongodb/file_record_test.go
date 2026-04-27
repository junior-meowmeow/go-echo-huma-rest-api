package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb/document"
)

func TestFileRecordRepository(t *testing.T) {
	db := setupMongoDatabase(t)
	ctx := context.Background()
	repo := mongodb.NewFileRecordRepository(db)
	coll := repo.Collection

	mockTime := time.Date(2025, 10, 25, 12, 0, 0, 0, time.UTC)

	t.Run("CreateFileRecord", func(t *testing.T) {
		cleanCollection(t, coll)

		t.Run("Should create file record successfully", func(t *testing.T) {
			input := &entity.FileRecord{
				FileName:    "test_photo.jpg",
				S3Key:       "uploads/random_name.jpg",
				ContentType: "image/jpeg",
				Size:        1048576,
				CreatedAt:   mockTime,
			}

			insertedID, err := repo.CreateFileRecord(ctx, input)

			require.NoError(t, err)
			assert.NotEmpty(t, insertedID)

			var doc document.FileRecordDocument
			oid, _ := document.StringToObjectID(insertedID)
			err = coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)

			require.NoError(t, err)
			assert.Equal(t, input.FileName, doc.FileName)
			assert.Equal(t, input.S3Key, doc.S3Key)
			assert.Equal(t, input.Size, doc.Size)
		})
	})

	t.Run("GetFileRecordByID", func(t *testing.T) {
		t.Run("Should return record when exists", func(t *testing.T) {
			cleanCollection(t, coll)

			rawOID := bson.NewObjectID()
			testDoc := document.FileRecordDocument{
				ID:          rawOID,
				FileName:    "report.pdf",
				S3Key:       "docs/report.pdf",
				ContentType: "application/pdf",
				CreatedAt:   mockTime,
			}
			_, err := coll.InsertOne(ctx, testDoc)
			require.NoError(t, err)

			idStr, _ := document.IDToString(rawOID)

			record, err := repo.GetFileRecordByID(ctx, idStr)

			require.NoError(t, err)
			assert.Equal(t, idStr, record.ID)
			assert.Equal(t, "report.pdf", record.FileName)
			assert.Equal(t, "application/pdf", record.ContentType)
		})

		t.Run("Should return error for invalid ID format", func(t *testing.T) {
			_, err := repo.GetFileRecordByID(ctx, "invalid-hex-string")

			require.Error(t, err)
			require.ErrorContains(t, err, "invalid file record ID format")
		})

		t.Run("Should return ErrNotFound when file record does not exist", func(t *testing.T) {
			cleanCollection(t, coll)

			randomID := bson.NewObjectID()
			idStr, _ := document.IDToString(randomID)

			_, err := repo.GetFileRecordByID(ctx, idStr)

			require.Error(t, err)
			require.ErrorIs(t, err, entity.ErrNotFound)
		})
	})
}
