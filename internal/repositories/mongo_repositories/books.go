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

type BooksRepository interface {
	CreateBook(ctx context.Context, book *entities.Book) (string, error)
	GetBookByID(ctx context.Context, id string) (entities.Book, error)
	GetAllBooks(ctx context.Context) ([]entities.Book, error)
	GetBooksWithPagination(ctx context.Context, pageSize int64, pageNumber int64) ([]entities.Book, error)
}

type booksRepository struct {
	Collection *mongo.Collection
}

func NewBooksRepository(db *mongo.Database) *booksRepository {
	return &booksRepository{
		Collection: db.Collection("books"),
	}
}

func (r *booksRepository) CreateBook(ctx context.Context, book *entities.Book) (string, error) {
	document, err := documents.NewBookDocument(book)
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

func (r *booksRepository) GetBookByID(ctx context.Context, id string) (entities.Book, error) {
	var book entities.Book

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return book, fmt.Errorf("invalid book ID format")
	}

	var document documents.BookDocument
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

func (r *booksRepository) GetAllBooks(ctx context.Context) ([]entities.Book, error) {
	var documents []documents.BookDocument

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

	books := make([]entities.Book, len(documents))
	for i, document := range documents {
		books[i] = document.ToEntity()
	}

	return books, nil
}

func (r *booksRepository) GetBooksWithPagination(ctx context.Context, pageSize int64, pageNumber int64) ([]entities.Book, error) {
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

	var documents []documents.BookDocument
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, fmt.Errorf("failed to decode book documents: %w", err)
	}

	books := make([]entities.Book, len(documents))
	for i, document := range documents {
		books[i] = document.ToEntity()
	}

	return books, nil
}
