//nolint:dupl // Documents are intended to follow a similar pattern.
package document

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
)

type UserDocument struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	Username string `bson:"username"`
	Password string `bson:"password"`
	Role     string `bson:"role"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

func NewUserDocument(user *entity.User) (UserDocument, error) {
	var userDocument UserDocument
	var err error

	oid, err := StringToObjectID(user.ID)
	if err != nil {
		return userDocument, fmt.Errorf("invalid user ID format: %w", err)
	}

	userDocument = UserDocument{
		ID:        oid,
		Username:  user.Username,
		Password:  user.Password,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userDocument, nil
}

func (doc *UserDocument) ToEntity() entity.User {
	return entity.User{
		ID:        doc.ID.Hex(),
		Username:  doc.Username,
		Password:  doc.Password,
		Role:      doc.Role,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}
