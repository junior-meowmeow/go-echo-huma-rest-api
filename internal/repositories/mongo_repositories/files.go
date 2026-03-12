package mongo_repositories

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories/mongo_repositories/documents"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FileMetadataRepository interface {
	CreateFileMetadata(ctx context.Context, fileMetadata *entities.FileMetadata) (string, error)
	GetFileMetadataByID(ctx context.Context, fileID string) (entities.FileMetadata, error)
}

type fileMetadataRepository struct {
	Collection *mongo.Collection
}

func NewFileMetadataRepository(db *mongo.Database) *fileMetadataRepository {
	return &fileMetadataRepository{
		Collection: db.Collection("file_metadata"),
	}
}

func (r *fileMetadataRepository) CreateFileMetadata(ctx context.Context, fileMetadata *entities.FileMetadata) (string, error) {
	document, err := documents.NewFileMetadataDocument(fileMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to convert file metadata to document: %w", err)
	}

	result, err := r.Collection.InsertOne(ctx, document)
	if err != nil {
		return "", fmt.Errorf("failed to insert file metadata document: %w", err)
	}

	insertedID := result.InsertedID.(bson.ObjectID).Hex()

	return insertedID, nil
}

func (r *fileMetadataRepository) GetFileMetadataByID(ctx context.Context, fileID string) (entities.FileMetadata, error) {
	var fileMetadata entities.FileMetadata

	oid, err := bson.ObjectIDFromHex(fileID)
	if err != nil {
		return fileMetadata, fmt.Errorf("invalid file metadata ID format")
	}

	var document documents.FileMetadataDocument
	err = r.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&document)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fileMetadata, fmt.Errorf("file metadata not found")
		}
		return fileMetadata, err
	}

	fileMetadata = document.ToEntity()

	return fileMetadata, nil
}
