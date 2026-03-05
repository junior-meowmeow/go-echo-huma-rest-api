package mongo_repository

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FileMetadataRepository interface {
	SaveFileMetadata(ctx context.Context, file *entities.FileMetadata) (string, error)
	GetFileMetadataByID(ctx context.Context, fileID string) (*entities.FileMetadata, error)
}

type MongoFilesRepository struct {
	Collection *mongo.Collection
}

func NewMongoFilesRepository(db *mongo.Database) *MongoFilesRepository {
	return &MongoFilesRepository{
		Collection: db.Collection("file_metadata"),
	}
}

func (r *MongoFilesRepository) SaveFileMetadata(ctx context.Context, record *entities.FileMetadata) (string, error) {
	res, err := r.Collection.InsertOne(ctx, record)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(bson.ObjectID).Hex(), nil
}

func (r *MongoFilesRepository) GetFileMetadataByID(ctx context.Context, fileID string) (*entities.FileMetadata, error) {
	oid, err := bson.ObjectIDFromHex(fileID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format")
	}

	var result entities.FileMetadata
	err = r.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
