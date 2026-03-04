package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FileMetadataRepository interface {
	SaveFileMetadata(ctx context.Context, file *FileRecord) (string, error)
	GetFileMetadataByID(ctx context.Context, fileID string) (*FileRecord, error)
}

type MongoFilesRepository struct {
	Collection *mongo.Collection
}

func NewMongoFilesRepository(db *mongo.Database) *MongoFilesRepository {
	return &MongoFilesRepository{
		Collection: db.Collection("files"),
	}
}

type FileRecord struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	Filename    string `bson:"filename"`
	Size        int64  `bson:"size"`
	ContentType string `bson:"contentType"`
	S3Key       string `bson:"s3Key"`

	CreatedAt  time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
}

func (r *MongoFilesRepository) SaveFileMetadata(ctx context.Context, record *FileRecord) (string, error) {
	res, err := r.Collection.InsertOne(ctx, record)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(bson.ObjectID).Hex(), nil
}

func (r *MongoFilesRepository) GetFileMetadataByID(ctx context.Context, fileID string) (*FileRecord, error) {
	oid, err := bson.ObjectIDFromHex(fileID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format")
	}

	var result FileRecord
	err = r.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
