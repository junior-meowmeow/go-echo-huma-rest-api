package usecase

import (
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/storage"
)

type UseCases struct {
	Greeting GreetingUseCase
	Review   ReviewUseCase
	File     FileUseCase
	Book     BookUseCase
	BookPage BookPageUseCase
	Health   HealthUseCase
}

func NewUseCases(repositories *repository.Repositories, storages *storage.Storages) *UseCases {
	return &UseCases{
		Greeting: NewGreetingUseCase(),
		Review:   NewReviewUseCase(repositories.Review),
		File:     NewFileUseCase(repositories.FileRecord, storages.ObjectStorage),
		Book:     NewBookUseCase(repositories.Book),
		BookPage: NewBookPageUseCase(repositories.Book, repositories.BookPage),
		Health:   NewHealthUseCase(),
	}
}
