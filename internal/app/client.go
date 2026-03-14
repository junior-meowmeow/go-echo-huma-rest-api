package app

import (
	"fmt"
	"log"
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

	log.Printf("Created a new PetStore client connected to %s\n", serverURL)

	return petStoreClient, nil
}
