package usecases

import "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories"

type UseCases struct {
	Greeting  GreetingUseCase
	Reviews   ReviewsUseCase
	Files     FilesUseCase
	Books     BooksUseCase
	BookPages BookPagesUseCase
	Health    HealthUseCase
}

func NewUseCases(repositories *repositories.Repositories) *UseCases {
	return &UseCases{
		Greeting:  NewGreetingUseCase(),
		Reviews:   NewReviewsUseCase(repositories.Reviews),
		Files:     NewFilesUseCase(repositories.FileMetadata, repositories.ObjectStorage),
		Books:     NewBooksUseCase(repositories.Books),
		BookPages: NewBookPagesUseCase(repositories.Books, repositories.BookPages),
		Health:    NewHealthUseCase(),
	}
}
