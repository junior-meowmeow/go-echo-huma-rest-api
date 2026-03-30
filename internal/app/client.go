package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	petStoreClient "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/external/petstore/client"
)

func newPetStoreClient(serverURL string, timeout time.Duration) (*petStoreClient.ClientWithResponses, error) {
	httpClient := &http.Client{
		Timeout: timeout,
	}
	petStoreClient, err := petStoreClient.NewClientWithResponses(
		serverURL,
		petStoreClient.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create PetStore client: %w", err)
	}

	slog.Info("Created a new PetStore client", slog.String("Petstore URL", serverURL))

	return petStoreClient, nil
}
