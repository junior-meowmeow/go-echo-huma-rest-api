package port

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (string, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
}
