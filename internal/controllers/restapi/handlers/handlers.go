package handlers

import (
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecases"
)

type Handlers struct {
	Greeting  GreetingHandler
	Reviews   ReviewsHandler
	Files     FilesHandler
	Books     BooksHandler
	BookPages BookPagesHandler
	Health    HealthHandler
}

func NewHandlers(usecases *usecases.UseCases) *Handlers {
	return &Handlers{
		Greeting:  NewGreetingHandler(usecases.Greeting),
		Reviews:   NewReviewsHandler(usecases.Reviews),
		Files:     NewFilesHandler(usecases.Files),
		Books:     NewBooksHandler(usecases.Books),
		BookPages: NewBookPagesHandler(usecases.BookPages),
		Health:    NewHealthHandler(usecases.Health),
	}
}
