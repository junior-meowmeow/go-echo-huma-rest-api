package usecase

import (
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository"
)

type UseCases struct {
	Greeting GreetingUseCase
	Review   ReviewUseCase
	File     FileUseCase
	Book     BookUseCase
	BookPage BookPageUseCase
	Health   HealthUseCase
}

func NewUseCases(repositories *repository.Repositories) *UseCases {
	return &UseCases{
		Greeting: NewGreetingUseCase(),
		Review:   NewReviewUseCase(repositories.Review),
		File:     NewFileUseCase(repositories.FileRecord, repositories.ObjectStorage),
		Book:     NewBookUseCase(repositories.Book),
		BookPage: NewBookPageUseCase(repositories.Book, repositories.BookPage),
		Health:   NewHealthUseCase(),
	}
}
