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

//revive:disable:unexported-return // Intentionally returns an unexported struct to enforce dependency on the interface in other layers.
func NewUserRepository(db *mongo.Database) *userRepository {
	return &userRepository{
		Collection: db.Collection("users"),
	}
}

//revive:enable:unexported-return

func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) (string, error) {
	doc, err := document.NewUserDocument(user)
	if err != nil {
		return "", fmt.Errorf("failed to convert user to document: %w", err)
	}

	result, err := r.Collection.InsertOne(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("failed to insert user document: %w", err)
	}

	insertedID, err := document.IDToString(result.InsertedID)
	if err != nil {
		return "", fmt.Errorf("failed to convert inserted id to string: %w", err)
	}

	return insertedID, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	var user entity.User

	var doc document.UserDocument
	err := r.Collection.FindOne(ctx, bson.M{"username": username}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, fmt.Errorf("failed to get user: %w: %w", entity.ErrNotFound, err)
		}
		return user, err
	}

	user = doc.ToEntity()

	return user, nil
}
