package document

import (
	"fmt"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserDocument struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	Username string `bson:"username"`
	Password string `bson:"password"`
	Role     string `bson:"role"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

func NewUserDocument(entity *entity.User) (UserDocument, error) {
	var userDocument UserDocument
	var err error

	var oid bson.ObjectID
	if entity.ID != "" {
		oid, err = bson.ObjectIDFromHex(entity.ID)
		if err != nil {
			return userDocument, fmt.Errorf("invalid user ID format: %w", err)
		}
	}

	userDocument = UserDocument{
		ID:        oid,
		Username:  entity.Username,
		Password:  entity.Password,
		Role:      entity.Role,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	return userDocument, nil
}

func (document *UserDocument) ToEntity() entity.User {
	return entity.User{
		ID:        document.ID.Hex(),
		Username:  document.Username,
		Password:  document.Password,
		Role:      document.Role,
		CreatedAt: document.CreatedAt,
		UpdatedAt: document.UpdatedAt,
	}
}
