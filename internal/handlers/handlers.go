package handlers

import "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories"

type Handlers struct {
	Greeting  GreetingHandler
	Reviews   ReviewsHandler
	Files     FilesHandler
	Books     BooksHandler
	BookPages BookPagesHandler
	Health    HealthHandler
}

func NewHandlers(repositories *repositories.Repositories) *Handlers {
	return &Handlers{
		Greeting:  NewGreetingHandler(),
		Reviews:   NewReviewsHandler(repositories.Reviews),
		Files:     NewFilesHandler(repositories.FileMetadata, repositories.ObjectStorage),
		Books:     NewBooksHandler(repositories.Books),
		BookPages: NewBookPagesHandler(repositories.Books, repositories.BookPages),
		Health:    NewHealthHandler(),
	}
}
