package handler

import (
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

type Handlers struct {
	Greeting GreetingHandler
	Review   ReviewHandler
	File     FileHandler
	Book     BookHandler
	BookPage BookPageHandler
	Health   HealthHandler
}

func NewHandlers(usecases *usecase.UseCases) *Handlers {
	return &Handlers{
		Greeting: NewGreetingHandler(usecases.Greeting),
		Review:   NewReviewHandler(usecases.Review),
		File:     NewFileHandler(usecases.File),
		Book:     NewBookHandler(usecases.Book),
		BookPage: NewBookPageHandler(usecases.BookPage),
		Health:   NewHealthHandler(usecases.Health),
	}
}
