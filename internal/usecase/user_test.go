package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/port/mocks"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/utility/auth"
	authmocks "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/utility/auth/mocks"
)

func TestUserUseCase(t *testing.T) {
	ctx := context.Background()

	t.Run("RegisterUser", func(t *testing.T) {
		t.Run("Should create user successfully", func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepository(t)
			mockToken := authmocks.NewMockTokenUtility(t)

			username := "new_user"
			password := "new_password"

			user := &entity.User{
				Username: username,
			}

			mockRepo.EXPECT().
				GetUserByUsername(ctx, username).
				Return(entity.User{}, entity.ErrNotFound)

			mockRepo.EXPECT().
				CreateUser(ctx, user).
				Run(func(_ context.Context, u *entity.User) {
					// Password should be hashed
					assert.NotEqual(t, password, u.Password)
					require.True(t, auth.CheckPasswordHash(password, u.Password))

					// Default role
					assert.Equal(t, "user", u.Role)

					// Timestamps
					assert.False(t, u.CreatedAt.IsZero())
					assert.False(t, u.UpdatedAt.IsZero())
					assert.WithinDuration(t, u.CreatedAt, u.UpdatedAt, time.Millisecond)
				}).
				Return("generated-id", nil)

			uc := usecase.NewUserUseCase(mockRepo, mockToken)

			id, err := uc.RegisterUser(ctx, user, password)

			require.NoError(t, err)
			assert.Equal(t, "generated-id", id)
		})

		t.Run("Should return error when username already exists", func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepository(t)
			mockToken := authmocks.NewMockTokenUtility(t)

			user := &entity.User{Username: "john"}

			mockRepo.EXPECT().
				GetUserByUsername(ctx, "john").
				Return(entity.User{}, nil)

			uc := usecase.NewUserUseCase(mockRepo, mockToken)

			id, err := uc.RegisterUser(ctx, user, "password")

			require.Error(t, err)
			assert.Empty(t, id)
			assert.ErrorIs(t, err, entity.ErrAlreadyExists)
		})

		t.Run("Should return error when checking username fails", func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepository(t)
			mockToken := authmocks.NewMockTokenUtility(t)

			mockErr := errors.New("db error")

			mockRepo.EXPECT().
				GetUserByUsername(ctx, "john").
				Return(entity.User{}, mockErr)

			uc := usecase.NewUserUseCase(mockRepo, mockToken)

			_, err := uc.RegisterUser(ctx, &entity.User{Username: "john"}, "password")

			require.Error(t, err)
			require.ErrorContains(t, err, "failed to check username")
			assert.ErrorIs(t, err, mockErr)
		})

		t.Run("Should return error when create user fails", func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepository(t)
			mockToken := authmocks.NewMockTokenUtility(t)

			mockErr := errors.New("insert error")

			mockRepo.EXPECT().
				GetUserByUsername(ctx, "john").
				Return(entity.User{}, entity.ErrNotFound)

			mockRepo.EXPECT().
				CreateUser(ctx, mock.Anything).
				Return("", mockErr)

			uc := usecase.NewUserUseCase(mockRepo, mockToken)

			_, err := uc.RegisterUser(ctx, &entity.User{Username: "john"}, "password")

			require.Error(t, err)
			require.ErrorContains(t, err, "failed to create user")
			assert.ErrorIs(t, err, mockErr)
		})
	})

	t.Run("LoginUser", func(t *testing.T) {
		t.Run("Should login successfully", func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepository(t)
			mockToken := authmocks.NewMockTokenUtility(t)

			hashedPassword, _ := auth.HashPassword("password")

			mockRepo.EXPECT().
				GetUserByUsername(ctx, "john").
				Return(entity.User{
					ID:       "1",
					Password: hashedPassword,
					Role:     "admin",
				}, nil)

			mockToken.EXPECT().
				GenerateToken("1", "admin").
				Return("token-123", nil)

			uc := usecase.NewUserUseCase(mockRepo, mockToken)

			token, err := uc.LoginUser(ctx, "john", "password")

			require.NoError(t, err)
			assert.Equal(t, "token-123", token)
		})

		t.Run("Should return error when user not found", func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepository(t)
			mockToken := authmocks.NewMockTokenUtility(t)

			mockErr := errors.New("db error")

			mockRepo.EXPECT().
				GetUserByUsername(ctx, "john").
				Return(entity.User{}, mockErr)

			uc := usecase.NewUserUseCase(mockRepo, mockToken)

			_, err := uc.LoginUser(ctx, "john", "password")

			require.Error(t, err)
			assert.ErrorIs(t, err, mockErr)
		})

		t.Run("Should return error when password is incorrect", func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepository(t)
			mockToken := authmocks.NewMockTokenUtility(t)

			hashedPassword, _ := auth.HashPassword("password")

			mockRepo.EXPECT().
				GetUserByUsername(ctx, "john").
				Return(entity.User{
					ID:       "1",
					Password: hashedPassword,
				}, nil)

			uc := usecase.NewUserUseCase(mockRepo, mockToken)

			_, err := uc.LoginUser(ctx, "john", "wrong")

			require.Error(t, err)
			assert.ErrorIs(t, err, entity.ErrInvalidCredentials)
		})

		t.Run("Should return error when token generation fails", func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepository(t)
			mockToken := authmocks.NewMockTokenUtility(t)

			hashedPassword, _ := auth.HashPassword("password")

			mockRepo.EXPECT().
				GetUserByUsername(ctx, "john").
				Return(entity.User{
					ID:       "1",
					Password: hashedPassword,
					Role:     "user",
				}, nil)

			mockErr := errors.New("token error")

			mockToken.EXPECT().
				GenerateToken("1", "user").
				Return("", mockErr)

			uc := usecase.NewUserUseCase(mockRepo, mockToken)

			_, err := uc.LoginUser(ctx, "john", "password")

			require.Error(t, err)
			require.ErrorContains(t, err, "failed to generate token")
			assert.ErrorIs(t, err, mockErr)
		})
	})
}
