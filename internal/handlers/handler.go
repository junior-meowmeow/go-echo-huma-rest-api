package handlers

import "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories"

type Handler struct {
	*repositories.Repositories
}

func NewHandler(repos *repositories.Repositories) *Handler {
	return &Handler{
		Repositories: repos,
	}
}
