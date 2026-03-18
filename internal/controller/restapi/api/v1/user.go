package v1

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
)

func RegisterUserGroup(public huma.API, protected huma.API, h *handler.Handlers) {
	publicGroup := huma.NewGroup(public, "/users")
	protectedGroup := huma.NewGroup(protected, "/users")

	RegisterUserRoutes(publicGroup, protectedGroup, h)
}

func RegisterUserRoutes(public huma.API, protected huma.API, h *handler.Handlers) {
	huma.Register(public, huma.Operation{
		OperationID:   "register-user",
		Method:        http.MethodPost,
		Path:          "/register",
		Summary:       "Register a new user",
		Description:   "Create a new user account with a username and password.",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusCreated,
	}, h.User.RegisterUser)

	huma.Register(public, huma.Operation{
		OperationID: "login-user",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "Log in a user",
		Description: "Authenticate a user and return a JWT Bearer token.",
		Tags:        []string{"Users"},
	}, h.User.LoginUser)
}
