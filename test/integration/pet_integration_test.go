package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PetSuite struct {
	IntegrationTestSuite
}

func TestPetSuite(t *testing.T) {
	suite.Run(t, new(PetSuite))
}

type petCategorySchema struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type petSchema struct {
	ID        int64             `json:"id"`
	Name      string            `json:"name"`
	Category  petCategorySchema `json:"category"`
	PhotoURLs []string          `json:"photoUrls"`
	Status    string            `json:"status"`
	Tags      []string          `json:"tags"`
}

type getAvailablePetsResponse struct {
	Data []petSchema `json:"data"`
}

func (s *PetSuite) getAvailablePets() *httptest.ResponseRecorder {
	s.T().Helper()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/pets", nil)
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *PetSuite) getPetByID(id int64) *httptest.ResponseRecorder {
	s.T().Helper()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/pets/"+intToStr(id), nil)
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *PetSuite) getPetByIDRaw(id string) *httptest.ResponseRecorder {
	s.T().Helper()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/pets/"+id, nil)
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *PetSuite) decodeGetAvailablePetsResponse(w *httptest.ResponseRecorder) getAvailablePetsResponse {
	s.T().Helper()

	var resp getAvailablePetsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err, "failed to decode GET /v1/pets response")
	return resp
}

func (s *PetSuite) decodeGetPetByIDResponse(w *httptest.ResponseRecorder) petSchema {
	s.T().Helper()

	var resp petSchema
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err, "failed to decode GET /v1/pets/{id} response")
	return resp
}

func (s *PetSuite) setMockPetHandler(fn http.HandlerFunc) {
	s.T().Helper()

	original := s.MockPetServer.Config.Handler
	s.MockPetServer.Config.Handler = fn
	s.T().Cleanup(func() {
		s.MockPetServer.Config.Handler = original
	})
}

func (s *PetSuite) TestGetAvailablePets_ReturnsHTTP200() {
	s.setMockPetHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))

	w := s.getAvailablePets()

	s.Equal(http.StatusOK, w.Code)
}

func (s *PetSuite) TestGetAvailablePets_EmptyList_ReturnsEmptyData() {
	s.setMockPetHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))

	w := s.getAvailablePets()
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetAvailablePetsResponse(w)
	s.NotNil(resp.Data)
	s.Empty(resp.Data)
}

func (s *PetSuite) TestGetAvailablePets_ForwardsStatusQueryToUpstream() {
	var capturedQuery string

	s.setMockPetHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.Query().Get("status")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))

	w := s.getAvailablePets()
	s.Require().Equal(http.StatusOK, w.Code)

	s.Equal("available", capturedQuery, "upstream must receive status=available")
}

func (s *PetSuite) TestGetAvailablePets_ReturnsCorrectFields() {
	mockPets := []map[string]any{
		{
			"id":        1,
			"name":      "Buddy",
			"status":    "available",
			"photoUrls": []string{"http://example.com/buddy.jpg"},
			"category":  map[string]any{"id": 10, "name": "Dogs"},
			"tags": []map[string]any{
				{"id": 1, "name": "friendly"},
				{"id": 2, "name": "trained"},
			},
		},
	}

	s.setMockPetHandler(jsonHandler(mockPets))

	w := s.getAvailablePets()
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetAvailablePetsResponse(w)
	s.Require().Len(resp.Data, 1)

	pet := resp.Data[0]
	s.Equal(int64(1), pet.ID)
	s.Equal("Buddy", pet.Name)
	s.Equal("available", pet.Status)
	s.Equal(int64(10), pet.Category.ID)
	s.Equal("Dogs", pet.Category.Name)
	s.Equal([]string{"http://example.com/buddy.jpg"}, pet.PhotoURLs)
	s.Equal([]string{"friendly", "trained"}, pet.Tags)
}

func (s *PetSuite) TestGetAvailablePets_MultiplePets_AllReturned() {
	mockPets := []map[string]any{
		{"id": 1, "name": "Dog", "status": "available", "photoUrls": []string{}},
		{"id": 2, "name": "Cat", "status": "available", "photoUrls": []string{}},
		{"id": 3, "name": "Mouse", "status": "available", "photoUrls": []string{}},
	}

	s.setMockPetHandler(jsonHandler(mockPets))

	w := s.getAvailablePets()
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetAvailablePetsResponse(w)
	s.Len(resp.Data, 3)
}

func (s *PetSuite) TestGetAvailablePets_UpstreamError_ReturnsServerError() {
	s.setMockPetHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	w := s.getAvailablePets()

	s.GreaterOrEqual(w.Code, 500, "upstream error should propagate as a 5xx")
}

func (s *PetSuite) TestGetAvailablePets_PetWithNoOptionalFields_Handled() {
	mockPets := []map[string]any{
		{"id": 42, "name": "Ghost", "photoUrls": []string{}, "status": "available"},
	}

	s.setMockPetHandler(jsonHandler(mockPets))

	w := s.getAvailablePets()
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetAvailablePetsResponse(w)
	s.Require().Len(resp.Data, 1)
	s.Equal(int64(42), resp.Data[0].ID)
	s.Equal("Ghost", resp.Data[0].Name)
}

func (s *PetSuite) TestGetPetByID_ReturnsHTTP200() {
	mockPet := map[string]any{
		"id": 123, "name": "TestPet", "status": "available", "photoUrls": []string{},
	}

	s.setMockPetHandler(jsonHandler(mockPet))

	w := s.getPetByID(123)

	s.Equal(http.StatusOK, w.Code)
}

func (s *PetSuite) TestGetPetByID_ReturnsCorrectPet() {
	mockPet := map[string]any{
		"id":        99,
		"name":      "Whiskers",
		"status":    "available",
		"photoUrls": []string{"http://example.com/whiskers.png"},
		"category":  map[string]any{"id": 5, "name": "Cats"},
		"tags": []map[string]any{
			{"id": 1, "name": "fluffy"},
		},
	}

	s.setMockPetHandler(jsonHandler(mockPet))

	w := s.getPetByID(99)
	s.Require().Equal(http.StatusOK, w.Code)

	pet := s.decodeGetPetByIDResponse(w)
	s.Equal(int64(99), pet.ID)
	s.Equal("Whiskers", pet.Name)
	s.Equal("available", pet.Status)
	s.Equal(int64(5), pet.Category.ID)
	s.Equal("Cats", pet.Category.Name)
	s.Equal([]string{"http://example.com/whiskers.png"}, pet.PhotoURLs)
	s.Equal([]string{"fluffy"}, pet.Tags)
}

func (s *PetSuite) TestGetPetByID_ForwardsIDToUpstream() {
	var capturedPath string

	s.setMockPetHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":456,"name":"Rex","photoUrls":[],"status":"available"}`))
	}))

	w := s.getPetByID(456)
	s.Require().Equal(http.StatusOK, w.Code)

	s.Equal("/pet/456", capturedPath, "upstream should receive the correct pet ID in the path")
}

func (s *PetSuite) TestGetPetByID_NotFound_Returns404() {
	s.setMockPetHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	w := s.getPetByID(999)

	s.Equal(http.StatusNotFound, w.Code)
}

func (s *PetSuite) TestGetPetByID_UpstreamError_ReturnsServerError() {
	s.setMockPetHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	w := s.getPetByID(1)

	s.GreaterOrEqual(w.Code, 500)
}

func (s *PetSuite) TestGetPetByID_InvalidIDFormat_Returns422() {
	cases := []struct {
		name string
		id   string
	}{
		{name: "Non-numeric", id: "abc"},
		{name: "Float", id: "1.5"},
		{name: "Empty-ish (spaces)", id: "   "},
		{name: "Special characters", id: "!@#"},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			w := s.getPetByIDRaw(tc.id)
			s.Equal(http.StatusUnprocessableEntity, w.Code, "case: %q", tc.name)
		})
	}
}

func (s *PetSuite) TestGetPetByID_PetWithNoOptionalFields_Handled() {
	s.setMockPetHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":7,"name":"Minimal","photoUrls":[],"status":"available"}`))
	}))

	w := s.getPetByID(7)
	s.Require().Equal(http.StatusOK, w.Code)

	pet := s.decodeGetPetByIDResponse(w)
	s.Equal(int64(7), pet.ID)
	s.Equal("Minimal", pet.Name)
	s.Empty(pet.Tags)
}

func jsonHandler(v any) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(v)
	}
}

func intToStr(n int64) string {
	return strconv.FormatInt(n, 10)
}
