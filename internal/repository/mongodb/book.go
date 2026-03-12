package mongodb

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository/mongodb/document"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BookRepository interface {
	CreateBook(ctx context.Context, book *entity.Book) (string, error)
	GetBookByID(ctx context.Context, id string) (entity.Book, error)
	GetAllBooks(ctx context.Context) ([]entity.Book, error)
	GetBooksWithPagination(ctx context.Context, pageSize int64, pageNumber int64) ([]entity.Book, error)
}

type bookRepository struct {
	Collection *mongo.Collection
}

func NewBookRepository(db *mongo.Database) *bookRepository {
	return &bookRepository{
		Collection: db.Collection("books"),
	}
}

func (r *bookRepository) CreateBook(ctx context.Context, book *entity.Book) (string, error) {
	document, err := document.NewBookDocument(book)
	if err != nil {
		return "", fmt.Errorf("failed to convert book to document: %w", err)
	}

	result, err := r.Collection.InsertOne(ctx, document)
	if err != nil {
		return "", fmt.Errorf("failed to insert book document: %w", err)
	}

	insertedID := result.InsertedID.(bson.ObjectID).Hex()

	return insertedID, nil
}

func (r *bookRepository) GetBookByID(ctx context.Context, id string) (entity.Book, error) {
	var book entity.Book

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return book, fmt.Errorf("invalid book ID format")
	}

	var document document.BookDocument
	err = r.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&document)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return book, fmt.Errorf("book not found")
		}
		return book, err
	}

	book = document.ToEntity()

	return book, nil
}

func (r *bookRepository) GetAllBooks(ctx context.Context) ([]entity.Book, error) {
	var documents []document.BookDocument

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.Collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &documents); err != nil {
		return nil, fmt.Errorf("failed to decode book documents: %w", err)
	}

	books := make([]entity.Book, len(documents))
	for i, document := range documents {
		books[i] = document.ToEntity()
	}

	return books, nil
}

func (r *bookRepository) GetBooksWithPagination(ctx context.Context, pageSize int64, pageNumber int64) ([]entity.Book, error) {
	skip := max((pageNumber-1)*pageSize, 0)

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(skip).
		SetLimit(pageSize)

	cursor, err := r.Collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []document.BookDocument
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, fmt.Errorf("failed to decode book documents: %w", err)
	}

	books := make([]entity.Book, len(documents))
	for i, document := range documents {
		books[i] = document.ToEntity()
	}

	return books, nil
}
