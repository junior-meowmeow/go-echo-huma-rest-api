package handlers

import "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository"

type Handler struct {
	*repository.Repositories
}

func NewHandler(repos *repository.Repositories) *Handler {
	return &Handler{
		Repositories: repos,
	}
}
