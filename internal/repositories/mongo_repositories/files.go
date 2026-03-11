package mongo_repositories

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FileMetadataRepository interface {
	SaveFileMetadata(ctx context.Context, file *entities.FileMetadata) (string, error)
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

func (r *fileMetadataRepository) SaveFileMetadata(ctx context.Context, record *entities.FileMetadata) (string, error) {
	result, err := r.Collection.InsertOne(ctx, record)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(bson.ObjectID).Hex(), nil
}

func (r *fileMetadataRepository) GetFileMetadataByID(ctx context.Context, fileID string) (entities.FileMetadata, error) {
	var fileMetadata entities.FileMetadata
	oid, err := bson.ObjectIDFromHex(fileID)
	if err != nil {
		return fileMetadata, fmt.Errorf("invalid ID format")
	}

	err = r.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&fileMetadata)
	if err != nil {
		return fileMetadata, err
	}
	return fileMetadata, nil
}
