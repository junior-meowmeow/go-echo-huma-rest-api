package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/port"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/utility/auth"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	RegisterUser(ctx context.Context, user *entity.User, password string) (string, error)
	LoginUser(ctx context.Context, username string, password string) (string, error)
}

type userUseCase struct {
	UserRepository port.UserRepository
	TokenUtility   auth.TokenUtility
}

func NewUserUseCase(userRepository port.UserRepository, tokenUtility auth.TokenUtility) *userUseCase {
	return &userUseCase{
		UserRepository: userRepository,
		TokenUtility:   tokenUtility,
	}
}

func (u *userUseCase) RegisterUser(ctx context.Context, user *entity.User, password string) (string, error) {
	_, err := u.UserRepository.GetUserByUsername(ctx, user.Username)
	if err == nil {
		return "", fmt.Errorf("username already taken: %w", entity.ErrAlreadyExists)
	}
	if !errors.Is(err, entity.ErrNotFound) {
		return "", fmt.Errorf("failed to check username: %w", err)
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt password: %w", err)
	}
	user.Password = string(hashedBytes)

	if user.Role == "" {
		user.Role = "user"
	}

	currentTime := time.Now()
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	id, err := u.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (u *userUseCase) LoginUser(ctx context.Context, username string, password string) (string, error) {
	user, err := u.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials: %w", entity.ErrInvalidCredentials)
	}

	token, err := u.TokenUtility.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}

	return token, nil
}
