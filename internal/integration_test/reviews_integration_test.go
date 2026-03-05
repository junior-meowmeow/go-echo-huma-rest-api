package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories/mongo_repositories"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ReviewsSuite struct {
	IntegrationTestSuite
}

func TestReviewsSuite(t *testing.T) {
	suite.Run(t, new(ReviewsSuite))
}

func (s *ReviewsSuite) SetupTest() {
	if repo, ok := s.Repos.Reviews.(*mongo_repositories.MongoReviewsRepository); ok {
		cleanCollection(s.T(), repo.Collection)
	} else {
		s.Fail("ReviewsRepository is not a MongoReviewsRepository, cannot clean collection")
	}
}

func (s *ReviewsSuite) TestCreateAndGetReview() {
	// 1. POST New Review
	requestBody := struct {
		Author  string `json:"author"`
		Rating  int    `json:"rating"`
		Message string `json:"message"`
	}{
		Author:  "Test",
		Rating:  5,
		Message: "Test Message",
	}

	jsonPayload, err := json.Marshal(requestBody)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/reviews", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	s.Router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		s.T().Logf("Response Status: %d", w.Code)
		s.T().Logf("Response Body: %s", w.Body.String())
	}

	s.Assert().Equal(http.StatusCreated, w.Code)

	// 2. GET Reviews
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/reviews", nil)
	s.Router.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusOK, w.Code)

	var responseBody []struct {
		ID        string `json:"id"`
		Author    string `json:"author"`
		Rating    int    `json:"rating"`
		Message   string `json:"message"`
		CreatedAt string `json:"createdAt"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	s.Require().NoError(err)

	s.Require().NotEmpty(responseBody)

	firstReview := responseBody[0]
	s.Equal("Test", firstReview.Author)
	s.Equal(5, firstReview.Rating)
	s.Equal("Test Message", firstReview.Message)
	s.NotEmpty(firstReview.ID)
	s.NotEmpty(firstReview.CreatedAt)
}

func (s *ReviewsSuite) TestCreateReviewWithInvalidRating() {
	// 1. POST Invalid Review (Rating > 5)
	requestBody := struct {
		Author  string `json:"author"`
		Rating  int    `json:"rating"`
		Message string `json:"message"`
	}{
		Author:  "Test",
		Rating:  6,
		Message: "This should fail",
	}

	jsonPayload, err := json.Marshal(requestBody)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/reviews", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	s.Router.ServeHTTP(w, req)

	s.Assert().NotEqual(http.StatusCreated, w.Code, "Should not create resource")
	s.Assert().Equal(http.StatusUnprocessableEntity, w.Code, "Should return status 422")

	// 2. GET Reviews (Verify that new review is not created)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/reviews", nil)
	s.Router.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusOK, w.Code)

	var responseBody []interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	s.Require().NoError(err)

	s.Assert().Len(responseBody, 0, "Database should be empty after failed creation")
}

func (s *ReviewsSuite) TestCreateReview_ValidationErrors() {
	type reviewBody struct {
		Author  string `json:"author"`
		Rating  int    `json:"rating"`
		Message string `json:"message"`
	}

	// Define the test cases
	tests := []struct {
		name       string
		body       reviewBody
		wantStatus int
	}{
		{
			name: "Rating Too High (6)",
			body: reviewBody{
				Author:  "ValidUser",
				Rating:  6, // Invalid > 5
				Message: "Should fail",
			},
			wantStatus: http.StatusUnprocessableEntity, // 422
		},
		{
			name: "Rating Too Low (0)",
			body: reviewBody{
				Author:  "ValidUser",
				Rating:  0, // Invalid < 1
				Message: "Should fail",
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Author Name Too Long (>10 chars)",
			body: reviewBody{
				Author:  "ThisNameIsWayTooLong", // Invalid > 10 chars
				Rating:  5,
				Message: "Should fail",
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Message Too Long (>100 chars)",
			body: reviewBody{
				Author: "ValidUser",
				Rating: 5,
				// Generate a string longer than 100 chars
				Message: string(make([]byte, 101)),
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	// Loop through test cases
	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Serialize Payload
			jsonPayload, err := json.Marshal(tt.body)
			s.Require().NoError(err)

			// Make Request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/v1/reviews", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")

			s.Router.ServeHTTP(w, req)

			// Assert Status Code
			s.Assert().Equal(tt.wantStatus, w.Code, "Expected failure status for case: %s", tt.name)

			// Verify that new review is not created
			// Access the collection directly to ensure count is 0.
			count, err := s.MongoDB.Collection("reviews").CountDocuments(context.Background(), bson.D{})
			s.Require().NoError(err)
			s.Assert().Equal(int64(0), count, "DB should be empty for case: %s", tt.name)
		})
	}
}
