package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BooksRepository interface {
	CreateBook(ctx context.Context, record *BookRecord) (string, error)
	GetBookByID(ctx context.Context, id string) (*BookRecord, error)
	GetAllBooks(ctx context.Context) ([]BookRecord, error)
	GetBooksWithPagination(ctx context.Context, pageSize int64, pageNumber int64) ([]BookRecord, error)
}

type MongoBooksRepository struct {
	Collection *mongo.Collection
}

func NewMongoBooksRepository(db *mongo.Database) *MongoBooksRepository {
	return &MongoBooksRepository{
		Collection: db.Collection("books"),
	}
}

type BookMetadata struct {
	Author string `bson:"author"`
	ISBN   string `bson:"isbn"`
	Genre  string `bson:"genre"`
}

type BookRecord struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	Name             string       `bson:"name"`
	Description      string       `bson:"description"`
	Metadata         BookMetadata `bson:"metadata"`
	CoverImageFileID string       `bson:"coverImageFileID"`

	CreatedAt  time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
}

func (r *MongoBooksRepository) CreateBook(ctx context.Context, record *BookRecord) (string, error) {
	res, err := r.Collection.InsertOne(ctx, record)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(bson.ObjectID).Hex(), nil
}

func (r *MongoBooksRepository) GetBookByID(ctx context.Context, id string) (*BookRecord, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID format")
	}

	var result BookRecord
	err = r.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("book not found")
		}
		return nil, err
	}

	return &result, nil
}

func (r *MongoBooksRepository) GetAllBooks(ctx context.Context) ([]BookRecord, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cur, err := r.Collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	results := make([]BookRecord, 0)
	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoBooksRepository) GetBooksWithPagination(ctx context.Context, pageSize int64, pageNumber int64) ([]BookRecord, error) {
	skip := max((pageNumber-1)*pageSize, 0)

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(skip).
		SetLimit(pageSize)

	cur, err := r.Collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	results := make([]BookRecord, 0, pageSize)
	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
