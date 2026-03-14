package usecase

import (
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/external"
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
	Pet      PetUseCase
}

func NewUseCases(repositories *repository.Repositories, storages *storage.Storages, services *external.ExternalServices) *UseCases {
	return &UseCases{
		Greeting: NewGreetingUseCase(),
		Review:   NewReviewUseCase(repositories.Review),
		File:     NewFileUseCase(repositories.FileRecord, storages.FileStorage),
		Book:     NewBookUseCase(repositories.Book),
		BookPage: NewBookPageUseCase(repositories.Book, repositories.BookPage),
		Health:   NewHealthUseCase(),
		Pet:      NewPetUseCase(services.PetService),
	}
}
