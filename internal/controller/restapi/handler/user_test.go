package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	usecasemocks "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase/mocks"
)

func TestUserHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("RegisterUser", func(t *testing.T) {
		t.Run("Should register user successfully", func(t *testing.T) {
			mockUC := usecasemocks.NewMockUserUseCase(t)

			req := &schema.RegisterUserRequest{}
			req.Body.Username = "john"
			req.Body.Password = "password"
			req.Body.Role = "admin"

			mockUC.EXPECT().
				RegisterUser(mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.Username == "john" && u.Role == "admin"
				}), "password").
				Return("user-id", nil)

			h := handler.NewUserHandler(mockUC)

			resp, err := h.RegisterUser(ctx, req)

			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, "user-id", resp.Body.ID)
		})

		t.Run("Should return 409 when username already exists", func(t *testing.T) {
			mockUC := usecasemocks.NewMockUserUseCase(t)

			req := &schema.RegisterUserRequest{}
			req.Body.Username = "john"
			req.Body.Password = "password"

			mockUC.EXPECT().
				RegisterUser(mock.Anything, mock.Anything, "password").
				Return("", entity.ErrAlreadyExists)

			h := handler.NewUserHandler(mockUC)

			resp, err := h.RegisterUser(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)

			require.ErrorContains(t, err, "Username is already taken")
		})

		t.Run("Should fallback to resolveError for other errors", func(t *testing.T) {
			mockUC := usecasemocks.NewMockUserUseCase(t)

			mockErr := errors.New("db error")

			req := &schema.RegisterUserRequest{}
			req.Body.Username = "john"
			req.Body.Password = "password"

			mockUC.EXPECT().
				RegisterUser(mock.Anything, mock.Anything, "password").
				Return("", mockErr)

			h := handler.NewUserHandler(mockUC)

			resp, err := h.RegisterUser(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)

			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})

	t.Run("LoginUser", func(t *testing.T) {
		t.Run("Should login successfully", func(t *testing.T) {
			mockUC := usecasemocks.NewMockUserUseCase(t)

			req := &schema.LoginUserRequest{}
			req.Body.Username = "john"
			req.Body.Password = "password"

			mockUC.EXPECT().
				LoginUser(mock.Anything, "john", "password").
				Return("jwt-token", nil)

			h := handler.NewUserHandler(mockUC)

			resp, err := h.LoginUser(ctx, req)

			require.NoError(t, err)
			require.NotNil(t, resp)

			require.Equal(t, "jwt-token", resp.Body.Token)
		})

		t.Run("Should return 401 for invalid credentials", func(t *testing.T) {
			mockUC := usecasemocks.NewMockUserUseCase(t)

			req := &schema.LoginUserRequest{}
			req.Body.Username = "john"
			req.Body.Password = "wrong"

			mockUC.EXPECT().
				LoginUser(mock.Anything, "john", "wrong").
				Return("", entity.ErrInvalidCredentials)

			h := handler.NewUserHandler(mockUC)

			resp, err := h.LoginUser(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)

			require.ErrorContains(t, err, "Invalid username or password")
		})

		t.Run("Should return 401 when user not found", func(t *testing.T) {
			mockUC := usecasemocks.NewMockUserUseCase(t)

			req := &schema.LoginUserRequest{}
			req.Body.Username = "john"
			req.Body.Password = "password"

			mockUC.EXPECT().
				LoginUser(mock.Anything, "john", "password").
				Return("", entity.ErrNotFound)

			h := handler.NewUserHandler(mockUC)

			resp, err := h.LoginUser(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)

			require.ErrorContains(t, err, "Invalid username or password")
		})

		t.Run("Should fallback to resolveError for other errors", func(t *testing.T) {
			mockUC := usecasemocks.NewMockUserUseCase(t)

			mockErr := errors.New("internal error")

			req := &schema.LoginUserRequest{}
			req.Body.Username = "john"
			req.Body.Password = "password"

			mockUC.EXPECT().
				LoginUser(mock.Anything, "john", "password").
				Return("", mockErr)

			h := handler.NewUserHandler(mockUC)

			resp, err := h.LoginUser(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)

			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})
}
