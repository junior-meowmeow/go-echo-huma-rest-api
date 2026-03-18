package handler

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
)

// Default error handler
func resolveError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, entity.ErrNotFound) {
		return huma.Error404NotFound("Resource not found", err)
	}

	if errors.Is(err, entity.ErrAlreadyExists) {
		return huma.Error409Conflict("Resource already exists", err)
	}

	if errors.Is(err, entity.ErrInvalidCredentials) {
		return huma.Error401Unauthorized("Invalid credentials", err)
	}

	return err
}
