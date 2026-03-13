package mongodb

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb/document"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type fileRecordRepository struct {
	Collection *mongo.Collection
}

func NewFileRecordRepository(db *mongo.Database) *fileRecordRepository {
	return &fileRecordRepository{
		Collection: db.Collection("filerecords"),
	}
}

func (r *fileRecordRepository) CreateFileRecord(ctx context.Context, fileRecord *entity.FileRecord) (string, error) {
	document, err := document.NewFileRecordDocument(fileRecord)
	if err != nil {
		return "", fmt.Errorf("failed to convert file record to document: %w", err)
	}

	result, err := r.Collection.InsertOne(ctx, document)
	if err != nil {
		return "", fmt.Errorf("failed to insert file record document: %w", err)
	}

	insertedID := result.InsertedID.(bson.ObjectID).Hex()

	return insertedID, nil
}

func (r *fileRecordRepository) GetFileRecordByID(ctx context.Context, fileID string) (entity.FileRecord, error) {
	var fileRecord entity.FileRecord

	oid, err := bson.ObjectIDFromHex(fileID)
	if err != nil {
		return fileRecord, fmt.Errorf("invalid file record ID format")
	}

	var document document.FileRecordDocument
	err = r.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&document)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fileRecord, fmt.Errorf("file record not found")
		}
		return fileRecord, err
	}

	fileRecord = document.ToEntity()

	return fileRecord, nil
}
