package middleware

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/util/auth"
)

func RequireToken(api huma.API, tokenUtility auth.TokenUtility) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		authHeader := ctx.Header("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "missing or invalid Authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		authContext, err := tokenUtility.ParseToken(tokenString)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid or expired token")
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
			huma.WriteErr(api, ctx, http.StatusForbidden, "insufficient permissions")
			return
		}

		next(ctx)
	}
}
