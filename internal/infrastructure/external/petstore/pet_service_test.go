package petstore_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/external/petstore"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/external/petstore/client"
)

func TestGetPetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("Should success", func(t *testing.T) {
		mockPet := client.Pet{
			Id:     ptr(int64(123)),
			Name:   "TestPet",
			Status: ptr(client.PetStatusAvailable),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/pet/123", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(mockPet)
		}))
		defer server.Close()

		apiClient, _ := client.NewClientWithResponses(server.URL)
		petService := petstore.NewPetService(apiClient)

		pet, err := petService.GetPetByID(ctx, 123)

		require.NoError(t, err)
		assert.Equal(t, int64(123), pet.ID)
		assert.Equal(t, "TestPet", pet.Name)
	})

	t.Run("Should handle network error correctly", func(t *testing.T) {
		apiClient, _ := client.NewClientWithResponses("http://invalid-address")
		petService := petstore.NewPetService(apiClient)

		_, err := petService.GetPetByID(ctx, 123)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "network error calling petstore api")
	})

	t.Run("Should handle not found error correctly", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		apiClient, _ := client.NewClientWithResponses(server.URL)
		petService := petstore.NewPetService(apiClient)

		pet, err := petService.GetPetByID(ctx, 999)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "pet not found")
		assert.Empty(t, pet.ID)
	})

	t.Run("Should handle unexpected code correctly", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer server.Close()

		apiClient, _ := client.NewClientWithResponses(server.URL)
		petService := petstore.NewPetService(apiClient)

		pet, err := petService.GetPetByID(ctx, 999)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code from petstore client")
		assert.Empty(t, pet.ID)
	})

	t.Run("Should handle empty response body correctly", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/pet/123", r.URL.Path)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		apiClient, _ := client.NewClientWithResponses(server.URL)
		petService := petstore.NewPetService(apiClient)

		pet, err := petService.GetPetByID(ctx, 123)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "received 200 OK but body was empty")
		assert.Empty(t, pet.ID)
	})
}

func TestGetPetsByStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("Should success", func(t *testing.T) {
		mockPets := []client.Pet{
			{Id: ptr(int64(1)), Name: "Dog", Status: ptr(client.PetStatusAvailable)},
			{Id: ptr(int64(2)), Name: "Cat", Status: ptr(client.PetStatusAvailable)},
			{Id: ptr(int64(3)), Name: "Mouse", Status: ptr(client.PetStatusAvailable)},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "available", r.URL.Query().Get("status"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(mockPets)
		}))
		defer server.Close()

		apiClient, _ := client.NewClientWithResponses(server.URL)
		petService := petstore.NewPetService(apiClient)

		pets, err := petService.GetPetsByStatus(ctx, entity.PetStatusAvailable)

		require.NoError(t, err)
		assert.Len(t, pets, 3)
		assert.Equal(t, "Dog", pets[0].Name)
		assert.Equal(t, "Cat", pets[1].Name)
		assert.Equal(t, "Mouse", pets[2].Name)
	})

	t.Run("Should handle network error correctly", func(t *testing.T) {
		apiClient, _ := client.NewClientWithResponses("http://invalid-address")
		petService := petstore.NewPetService(apiClient)

		_, err := petService.GetPetsByStatus(ctx, entity.PetStatusAvailable)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "network error calling petstore api")
	})

	t.Run("Should handle unexpected code correctly", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer server.Close()

		apiClient, _ := client.NewClientWithResponses(server.URL)
		petService := petstore.NewPetService(apiClient)

		_, err := petService.GetPetsByStatus(ctx, entity.PetStatusAvailable)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code from petstore client")
	})

	t.Run("Should handle empty response body correctly", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		apiClient, _ := client.NewClientWithResponses(server.URL)
		petService := petstore.NewPetService(apiClient)

		pets, err := petService.GetPetsByStatus(ctx, entity.PetStatusAvailable)

		require.NoError(t, err)
		assert.Empty(t, pets)
	})
}

// Helper to handle pointer types generated by oapi-codegen.
func ptr[T any](v T) *T {
	return &v
}
