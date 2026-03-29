package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb/document"
)

type fileRecordRepository struct {
	Collection *mongo.Collection
}

//revive:disable:unexported-return // Intentionally returns an unexported struct to enforce dependency on the interface in other layers.
func NewFileRecordRepository(db *mongo.Database) *fileRecordRepository {
	return &fileRecordRepository{
		Collection: db.Collection("filerecords"),
	}
}

//revive:enable:unexported-return

func (r *fileRecordRepository) CreateFileRecord(ctx context.Context, fileRecord *entity.FileRecord) (string, error) {
	doc, err := document.NewFileRecordDocument(fileRecord)
	if err != nil {
		return "", fmt.Errorf("failed to convert file record to document: %w", err)
	}

	result, err := r.Collection.InsertOne(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("failed to insert file record document: %w", err)
	}

	insertedID, err := document.IDToString(result.InsertedID)
	if err != nil {
		return "", fmt.Errorf("failed to convert inserted id to string: %w", err)
	}

	return insertedID, nil
}

func (r *fileRecordRepository) GetFileRecordByID(ctx context.Context, fileID string) (entity.FileRecord, error) {
	var fileRecord entity.FileRecord

	oid, err := document.StringToObjectID(fileID)
	if err != nil {
		return fileRecord, fmt.Errorf("invalid file record ID format: %w", err)
	}

	var doc document.FileRecordDocument
	err = r.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fileRecord, fmt.Errorf("failed to get file record: %w: %w", entity.ErrNotFound, err)
		}
		return fileRecord, err
	}

	fileRecord = doc.ToEntity()

	return fileRecord, nil
}
