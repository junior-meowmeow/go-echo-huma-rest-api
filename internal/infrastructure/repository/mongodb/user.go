package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb/document"
)

type userRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *userRepository {
	return &userRepository{
		Collection: db.Collection("users"),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) (string, error) {
	document, err := document.NewUserDocument(user)
	if err != nil {
		return "", fmt.Errorf("failed to convert user to document: %w", err)
	}

	result, err := r.Collection.InsertOne(ctx, document)
	if err != nil {
		return "", fmt.Errorf("failed to insert user document: %w", err)
	}

	insertedID := result.InsertedID.(bson.ObjectID).Hex()

	return insertedID, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	var user entity.User

	var document document.UserDocument
	err := r.Collection.FindOne(ctx, bson.M{"username": username}).Decode(&document)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, fmt.Errorf("failed to get user: %w: %w", entity.ErrNotFound, err)
		}
		return user, err
	}

	user = document.ToEntity()

	return user, nil
}
