package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository/mongodb/document"
)

func TestUserRepository(t *testing.T) {
	db := setupMongoDatabase(t)
	ctx := context.Background()
	repo := mongodb.NewUserRepository(db)
	coll := repo.Collection

	t.Run("CreateUser", func(t *testing.T) {
		cleanCollection(t, coll)

		t.Run("Should create user successfully", func(t *testing.T) {
			input := &entity.User{
				Username: "new_user",
				Password: "new_password",
				Role:     "admin",
			}

			insertedID, err := repo.CreateUser(ctx, input)

			require.NoError(t, err)
			assert.NotEmpty(t, insertedID)

			var doc document.UserDocument
			err = coll.FindOne(ctx, bson.M{"username": "new_user"}).Decode(&doc)

			require.NoError(t, err)
			assert.Equal(t, input.Username, doc.Username)
			assert.Equal(t, input.Password, doc.Password)
			assert.Equal(t, input.Role, doc.Role)
		})
	})

	t.Run("GetUserByUsername", func(t *testing.T) {
		t.Run("Should return user when exists", func(t *testing.T) {
			cleanCollection(t, coll)

			testDoc := document.UserDocument{
				Username: "test_user",
				Password: "test_password",
				Role:     "user",
			}
			_, err := coll.InsertOne(ctx, testDoc)
			require.NoError(t, err)

			user, err := repo.GetUserByUsername(ctx, "test_user")

			require.NoError(t, err)
			assert.Equal(t, testDoc.Username, user.Username)
			assert.Equal(t, testDoc.Role, user.Role)
		})

		t.Run("Should return ErrNotFound when user does not exist", func(t *testing.T) {
			cleanCollection(t, coll)

			_, err := repo.GetUserByUsername(ctx, "non_existent_user")

			require.Error(t, err)
			require.ErrorIs(t, err, entity.ErrNotFound)
		})
	})
}
