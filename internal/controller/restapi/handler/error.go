package handler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
)

// This is for local debugging only.
const isDebug = false

// ResolveError is default error handler.
func ResolveError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	reqID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		reqID = "unknown"
	}

	slog.ErrorContext(ctx, "HTTP_HANDLER_ERROR",
		slog.String("request_id", reqID),
		slog.Any("error", err),
	)

	return formatError(err)
}

func formatError(err error) error {
	var details []error
	if isDebug {
		details = append(details, err)
	}

	switch {
	case errors.Is(err, entity.ErrNotFound):
		return huma.Error404NotFound("Resource not found", details...)

	case errors.Is(err, entity.ErrAlreadyExists):
		return huma.Error409Conflict("Resource already exists", details...)

	case errors.Is(err, entity.ErrInvalidCredentials):
		return huma.Error401Unauthorized("Invalid credentials", details...)
	}

	if isDebug {
		return err
	}
	return huma.Error500InternalServerError("An unexpected internal error occurred")
}
