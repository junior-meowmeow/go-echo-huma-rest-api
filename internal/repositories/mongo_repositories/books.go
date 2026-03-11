package mongo_repositories

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BooksRepository interface {
	CreateBook(ctx context.Context, record *entities.Book) (string, error)
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

func (r *booksRepository) CreateBook(ctx context.Context, record *entities.Book) (string, error) {
	result, err := r.Collection.InsertOne(ctx, record)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(bson.ObjectID).Hex(), nil
}

func (r *booksRepository) GetBookByID(ctx context.Context, id string) (entities.Book, error) {
	var book entities.Book
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return book, fmt.Errorf("invalid book ID format")
	}

	err = r.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return book, fmt.Errorf("book not found")
		}
		return book, err
	}

	return book, nil
}

func (r *booksRepository) GetAllBooks(ctx context.Context) ([]entities.Book, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cur, err := r.Collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	books := make([]entities.Book, 0)
	if err := cur.All(ctx, &books); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *booksRepository) GetBooksWithPagination(ctx context.Context, pageSize int64, pageNumber int64) ([]entities.Book, error) {
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

	books := make([]entities.Book, 0, pageSize)
	if err := cur.All(ctx, &books); err != nil {
		return nil, err
	}
	return books, nil
}
