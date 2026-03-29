package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb/document"
)

type bookRepository struct {
	Collection *mongo.Collection
}

//revive:disable:unexported-return // Intentionally returns an unexported struct to enforce dependency on the interface in other layers.
func NewBookRepository(db *mongo.Database) *bookRepository {
	return &bookRepository{
		Collection: db.Collection("books"),
	}
}

//revive:enable:unexported-return

func (r *bookRepository) CreateBook(ctx context.Context, book *entity.Book) (string, error) {
	doc, err := document.NewBookDocument(book)
	if err != nil {
		return "", fmt.Errorf("failed to convert book to document: %w", err)
	}

	result, err := r.Collection.InsertOne(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("failed to insert book document: %w", err)
	}

	insertedID, err := document.IDToString(result.InsertedID)
	if err != nil {
		return "", fmt.Errorf("failed to convert inserted id to string: %w", err)
	}

	return insertedID, nil
}

func (r *bookRepository) GetBookByID(ctx context.Context, id string) (entity.Book, error) {
	var book entity.Book

	oid, err := document.StringToObjectID(id)
	if err != nil {
		return book, fmt.Errorf("invalid book ID format: %w", err)
	}

	var doc document.BookDocument
	err = r.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return book, fmt.Errorf("failed to get book: %w: %w", entity.ErrNotFound, err)
		}
		return book, err
	}

	book = doc.ToEntity()

	return book, nil
}

func (r *bookRepository) GetAllBooks(ctx context.Context) ([]entity.Book, error) {
	var docs []document.BookDocument

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.Collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode book documents: %w", err)
	}

	books := make([]entity.Book, len(docs))
	for i, doc := range docs {
		books[i] = doc.ToEntity()
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

	var docs []document.BookDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode book documents: %w", err)
	}

	books := make([]entity.Book, len(docs))
	for i, doc := range docs {
		books[i] = doc.ToEntity()
	}

	return books, nil
}
