package handler

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

type UserHandler interface {
	RegisterUser(ctx context.Context, request *schema.RegisterUserRequest) (*schema.RegisterUserResponse, error)
	LoginUser(ctx context.Context, request *schema.LoginUserRequest) (*schema.LoginUserResponse, error)
}

type userHandler struct {
	UserUseCase usecase.UserUseCase
}

//revive:disable:unexported-return // Intentionally returns an unexported struct to enforce dependency on the interface in other layers.
func NewUserHandler(userUseCase usecase.UserUseCase) *userHandler {
	return &userHandler{
		UserUseCase: userUseCase,
	}
}

//revive:enable:unexported-return

func (h *userHandler) RegisterUser(ctx context.Context, request *schema.RegisterUserRequest) (*schema.RegisterUserResponse, error) {
	user := &entity.User{
		Username: request.Body.Username,
		Role:     request.Body.Role,
	}

	id, err := h.UserUseCase.RegisterUser(ctx, user, request.Body.Password)
	if err != nil {
		if errors.Is(err, entity.ErrAlreadyExists) {
			return nil, huma.Error409Conflict("Username is already taken")
		}
		return nil, resolveError(err)
	}

	resp := &schema.RegisterUserResponse{}
	resp.Body.ID = id

	return resp, nil
}

func (h *userHandler) LoginUser(ctx context.Context, request *schema.LoginUserRequest) (*schema.LoginUserResponse, error) {
	token, err := h.UserUseCase.LoginUser(ctx, request.Body.Username, request.Body.Password)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) || errors.Is(err, entity.ErrNotFound) {
			return nil, huma.Error401Unauthorized("Invalid username or password")
		}
		return nil, resolveError(err)
	}

	resp := &schema.LoginUserResponse{}
	resp.Body.Token = token

	return resp, nil
}
