package handlers

import "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories"

type Handler struct {
	Greeting  GreetingHandler
	Reviews   ReviewsHandler
	Files     FilesHandler
	Books     BooksHandler
	BookPages BookPagesHandler
	Health    HealthHandler
}

func NewHandler(repos *repositories.Repositories) *Handler {
	return &Handler{
		Greeting:  NewGreetingHandler(),
		Reviews:   NewReviewsHandler(repos.Reviews),
		Files:     NewFilesHandler(repos.FileMetadata, repos.ObjectStorage),
		Books:     NewBooksHandler(repos.Books),
		BookPages: NewBookPagesHandler(repos.Books, repos.BookPages),
		Health:    NewHealthHandler(),
	}
}
