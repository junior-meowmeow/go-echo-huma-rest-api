package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BookPagesRepository interface {
	CreateBookPage(ctx context.Context, record *BookPageRecord) (string, error)
	GetBookPageByID(ctx context.Context, id string) (*BookPageRecord, error)
	GetBookPagesByBookID(ctx context.Context, bookID string) ([]BookPageRecord, error)
	GetBookpagesByBookIDWithPagination(ctx context.Context, bookID string, pageSize int64, pageNumber int64) ([]BookPageRecord, error)
	GetBookpagesByPageRange(ctx context.Context, bookID string, startPage int64, endPage int64) ([]BookPageRecord, error)
	GetBookpagesAroundPageNumber(ctx context.Context, bookID string, centerPage int64, offset int64) ([]BookPageRecord, error)
}

type MongoBookPagesRepository struct {
	Collection *mongo.Collection
}

func NewMongoBookPagesRepository(db *mongo.Database) *MongoBookPagesRepository {
	return &MongoBookPagesRepository{
		Collection: db.Collection("book_pages"),
	}
}

type BookPageMetadata struct {
	IsBookmarked bool   `bson:"isBookmarked"`
	Highlight    string `bson:"highlight"`
}

type BookPageRecord struct {
	ID     bson.ObjectID `bson:"_id,omitempty"`
	BookID bson.ObjectID `bson:"book_id,omitempty"`

	PageNumber          int64            `bson:"pageNumber"`
	Content             string           `bson:"content"`
	Metadata            BookPageMetadata `bson:"metadata"`
	AttachedImageFileID string           `bson:"attachedImageFileID,omitempty"`

	CreatedAt  time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
}

func (r *MongoBookPagesRepository) CreateBookPage(ctx context.Context, record *BookPageRecord) (string, error) {
	res, err := r.Collection.InsertOne(ctx, record)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(bson.ObjectID).Hex(), nil
}

func (r *MongoBookPagesRepository) GetBookPageByID(ctx context.Context, id string) (*BookPageRecord, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid book page ID format")
	}

	var result BookPageRecord
	err = r.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("book page not found")
		}
		return nil, err
	}
	return &result, nil
}

func (r *MongoBookPagesRepository) GetBookPagesByBookID(ctx context.Context, bookID string) ([]BookPageRecord, error) {
	b_oid, err := bson.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book page ID format")
	}

	filter := bson.D{{Key: "book_id", Value: b_oid}}

	opts := options.Find().
		SetSort(bson.D{{Key: "pageNumber", Value: 1}})

	cur, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	results := make([]BookPageRecord, 0)
	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoBookPagesRepository) GetBookpagesByBookIDWithPagination(ctx context.Context, bookID string, pageSize int64, pageNumber int64) ([]BookPageRecord, error) {
	b_oid, err := bson.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book page ID format")
	}

	skip := max((pageNumber-1)*pageSize, 0)
	filter := bson.D{{Key: "book_id", Value: b_oid}}

	opts := options.Find().
		SetSort(bson.D{{Key: "pageNumber", Value: 1}}).
		SetSkip(skip).
		SetLimit(pageSize)

	cur, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	results := make([]BookPageRecord, 0, pageSize)
	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoBookPagesRepository) GetBookpagesByPageRange(ctx context.Context, bookID string, startPage int64, endPage int64) ([]BookPageRecord, error) {
	b_oid, err := bson.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book page ID format")
	}

	filter := bson.D{
		{Key: "book_id", Value: b_oid},
		{Key: "pageNumber", Value: bson.D{
			{Key: "$gte", Value: startPage},
			{Key: "$lte", Value: endPage},
		}},
	}

	opts := options.Find().SetSort(bson.D{{Key: "pageNumber", Value: 1}})

	cur, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	results := make([]BookPageRecord, 0)
	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoBookPagesRepository) GetBookpagesAroundPageNumber(ctx context.Context, bookID string, centerPage int64, offset int64) ([]BookPageRecord, error) {
	e_oid, err := bson.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID format")
	}

	// Fetch "Past + Center" (<= page number)
	filterPast := bson.D{
		{Key: "book_id", Value: e_oid},
		{Key: "pageNumber", Value: bson.D{{Key: "$lte", Value: centerPage}}},
	}
	optsPast := options.Find().
		SetSort(bson.D{{Key: "pageNumber", Value: -1}}).
		SetLimit(offset + 1)

	cursorPast, err := r.Collection.Find(ctx, filterPast, optsPast)
	if err != nil {
		return nil, err
	}
	var past []BookPageRecord
	if err := cursorPast.All(ctx, &past); err != nil {
		return nil, err
	}

	// Fetch "Future" (> page number)
	var future []BookPageRecord
	if offset > 0 {
		filterFuture := bson.D{
			{Key: "book_id", Value: e_oid},
			{Key: "pageNumber", Value: bson.D{{Key: "$gt", Value: centerPage}}},
		}
		optsFuture := options.Find().
			SetSort(bson.D{{Key: "pageNumber", Value: 1}}).
			SetLimit(offset)

		cursorFuture, err := r.Collection.Find(ctx, filterFuture, optsFuture)
		if err != nil {
			return nil, err
		}
		if err := cursorFuture.All(ctx, &future); err != nil {
			return nil, err
		}
	}

	// Merge and Sort Data
	result := make([]BookPageRecord, 0, len(past)+len(future))
	for i := len(past) - 1; i >= 0; i-- {
		result = append(result, past[i])
	}
	result = append(result, future...)

	return result, nil
}
