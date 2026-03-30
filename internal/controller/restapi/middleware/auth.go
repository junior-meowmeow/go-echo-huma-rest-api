package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/utility/auth"
)

func RequireToken(api huma.API, tokenUtility auth.TokenUtility) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		authHeader := ctx.Header("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			err := huma.WriteErr(api, ctx, http.StatusUnauthorized, "missing or invalid Authorization header")
			if err != nil {
				slog.Error("Failed to write error response", slog.Any("error", err))
			}
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		authContext, err := tokenUtility.ParseToken(tokenString)
		if err != nil {
			err := huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid or expired token")
			if err != nil {
				slog.Error("Failed to write error response", slog.Any("error", err))
			}
			return
		}

		ctx = huma.WithValue(ctx, auth.AuthContextKey, authContext)
		next(ctx)
	}
}

func RequireAdminRole(api huma.API) func(huma.Context, func(huma.Context)) {
	return requireRole(api, "admin")
}

func requireRole(api huma.API, role string) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		authContext, ok := ctx.Context().Value(auth.AuthContextKey).(auth.AuthContext)

		if !ok || authContext.Role != role {
			err := huma.WriteErr(api, ctx, http.StatusForbidden, "insufficient permissions")
			if err != nil {
				slog.Error("Failed to write error response", slog.Any("error", err))
			}
			return
		}

		next(ctx)
	}
}
