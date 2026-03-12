package mongo_repositories

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories/mongo_repositories/documents"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BookPagesRepository interface {
	CreateBookPage(ctx context.Context, bookPage *entities.BookPage) (string, error)
	GetBookPageByID(ctx context.Context, id string) (entities.BookPage, error)
	GetBookPagesByBookID(ctx context.Context, bookID string) ([]entities.BookPage, error)
	GetBookpagesByBookIDWithPagination(ctx context.Context, bookID string, pageSize int64, pageNumber int64) ([]entities.BookPage, error)
	GetBookpagesByPageRange(ctx context.Context, bookID string, startPage int64, endPage int64) ([]entities.BookPage, error)
	GetBookpagesAroundPageNumber(ctx context.Context, bookID string, centerPage int64, offset int64) ([]entities.BookPage, error)
}

type bookPagesRepository struct {
	Collection *mongo.Collection
}

func NewBookPagesRepository(db *mongo.Database) *bookPagesRepository {
	return &bookPagesRepository{
		Collection: db.Collection("book_pages"),
	}
}

func (r *bookPagesRepository) CreateBookPage(ctx context.Context, bookPage *entities.BookPage) (string, error) {
	document, err := documents.NewBookPageDocument(bookPage)
	if err != nil {
		return "", fmt.Errorf("failed to convert book page to document: %w", err)
	}

	result, err := r.Collection.InsertOne(ctx, document)
	if err != nil {
		return "", fmt.Errorf("failed to insert book page document: %w", err)
	}

	insertedID := result.InsertedID.(bson.ObjectID).Hex()

	return insertedID, nil
}

func (r *bookPagesRepository) GetBookPageByID(ctx context.Context, id string) (entities.BookPage, error) {
	var bookPage entities.BookPage

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bookPage, fmt.Errorf("invalid book page ID format")
	}

	var document documents.BookPageDocument
	err = r.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&document)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return bookPage, fmt.Errorf("book page not found")
		}
		return bookPage, err
	}

	bookPage = document.ToEntity()

	return bookPage, nil
}

func (r *bookPagesRepository) GetBookPagesByBookID(ctx context.Context, bookID string) ([]entities.BookPage, error) {
	bookOID, err := bson.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book page ID format")
	}

	var documents []documents.BookPageDocument

	filter := bson.D{{Key: "book_id", Value: bookOID}}

	opts := options.Find().
		SetSort(bson.D{{Key: "pageNumber", Value: 1}})

	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &documents); err != nil {
		return nil, fmt.Errorf("failed to decode book page documents: %w", err)
	}

	bookPages := make([]entities.BookPage, len(documents))
	for i, document := range documents {
		bookPages[i] = document.ToEntity()
	}

	return bookPages, nil
}

func (r *bookPagesRepository) GetBookpagesByBookIDWithPagination(ctx context.Context, bookID string, pageSize int64, pageNumber int64) ([]entities.BookPage, error) {
	bookOID, err := bson.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book page ID format")
	}

	skip := max((pageNumber-1)*pageSize, 0)
	filter := bson.D{{Key: "book_id", Value: bookOID}}

	opts := options.Find().
		SetSort(bson.D{{Key: "pageNumber", Value: 1}}).
		SetSkip(skip).
		SetLimit(pageSize)

	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	bookPages := make([]entities.BookPage, 0, pageSize)
	if err := cursor.All(ctx, &bookPages); err != nil {
		return nil, err
	}
	return bookPages, nil
}

func (r *bookPagesRepository) GetBookpagesByPageRange(ctx context.Context, bookID string, startPage int64, endPage int64) ([]entities.BookPage, error) {
	bookOID, err := bson.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book page ID format")
	}

	filter := bson.D{
		{Key: "book_id", Value: bookOID},
		{Key: "pageNumber", Value: bson.D{
			{Key: "$gte", Value: startPage},
			{Key: "$lte", Value: endPage},
		}},
	}

	opts := options.Find().SetSort(bson.D{{Key: "pageNumber", Value: 1}})

	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	bookPages := make([]entities.BookPage, 0)
	if err := cursor.All(ctx, &bookPages); err != nil {
		return nil, err
	}
	return bookPages, nil
}

func (r *bookPagesRepository) GetBookpagesAroundPageNumber(ctx context.Context, bookID string, centerPage int64, offset int64) ([]entities.BookPage, error) {
	bookOID, err := bson.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID format")
	}

	// Fetch "Past + Center" (<= page number)
	filterPast := bson.D{
		{Key: "book_id", Value: bookOID},
		{Key: "pageNumber", Value: bson.D{{Key: "$lte", Value: centerPage}}},
	}
	optsPast := options.Find().
		SetSort(bson.D{{Key: "pageNumber", Value: -1}}).
		SetLimit(offset + 1)

	cursorPast, err := r.Collection.Find(ctx, filterPast, optsPast)
	if err != nil {
		return nil, err
	}
	var past []entities.BookPage
	if err := cursorPast.All(ctx, &past); err != nil {
		return nil, err
	}

	// Fetch "Future" (> page number)
	var future []entities.BookPage
	if offset > 0 {
		filterFuture := bson.D{
			{Key: "book_id", Value: bookOID},
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
	bookPages := make([]entities.BookPage, 0, len(past)+len(future))
	for i := len(past) - 1; i >= 0; i-- {
		bookPages = append(bookPages, past[i])
	}
	bookPages = append(bookPages, future...)

	return bookPages, nil
}
