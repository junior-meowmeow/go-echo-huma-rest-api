package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/test/testhelper"
)

type ReviewSuite struct {
	IntegrationTestSuite
}

func TestReviewSuite(t *testing.T) {
	suite.Run(t, new(ReviewSuite))
}

func (s *ReviewSuite) SetupTest() {
	testhelper.CleanMongoCollection(s.T(), s.MongoDB.Collection("reviews"))
}

type createReviewBody struct {
	Author  string `json:"author"`
	Rating  int    `json:"rating"`
	Message string `json:"message"`
}

type reviewSchema struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Rating    int       `json:"rating"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type getReviewsResponse struct {
	Data []reviewSchema `json:"data"`
}

func (s *ReviewSuite) postReview(body createReviewBody) *httptest.ResponseRecorder {
	s.T().Helper()

	payload, err := json.Marshal(body)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/reviews", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *ReviewSuite) getReviews() *httptest.ResponseRecorder {
	s.T().Helper()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/reviews", nil)
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *ReviewSuite) mustPostReview(body createReviewBody) {
	s.T().Helper()

	w := s.postReview(body)
	if w.Code != http.StatusCreated {
		s.T().Logf("mustPostReview failed - status: %d, body: %s", w.Code, w.Body.String())
	}
	s.Require().Equal(http.StatusCreated, w.Code)
}

func (s *ReviewSuite) decodeGetReviewsResponse(w *httptest.ResponseRecorder) getReviewsResponse {
	s.T().Helper()

	var resp getReviewsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err, "failed to decode GET /v1/reviews response")
	return resp
}

func (s *ReviewSuite) dbReviewCount() int64 {
	s.T().Helper()

	count, err := s.MongoDB.Collection("reviews").CountDocuments(context.Background(), bson.D{})
	s.Require().NoError(err)
	return count
}

func validReview() createReviewBody {
	return createReviewBody{
		Author:  "Alice",
		Rating:  5,
		Message: "Great service!",
	}
}

func (s *ReviewSuite) TestPostReview_CreatesReview_ReturnsHTTP201() {
	w := s.postReview(validReview())

	s.Equal(http.StatusCreated, w.Code)
}

func (s *ReviewSuite) TestPostReview_CreatesReview_PersistsToMongoDB() {
	body := validReview()
	s.mustPostReview(body)

	s.Equal(int64(1), s.dbReviewCount())

	var doc bson.M
	err := s.MongoDB.Collection("reviews").FindOne(context.Background(), bson.D{}).Decode(&doc)
	s.Require().NoError(err)

	s.Equal(body.Author, doc["author"])
	s.Equal(int32(body.Rating), doc["rating"])
	s.Equal(body.Message, doc["message"])
	s.NotEmpty(doc["_id"], "document should have a MongoDB _id")
	s.NotEmpty(doc["createdAt"], "document should have a createdAt timestamp")
}

func (s *ReviewSuite) TestPostReview_MultipleReviews_AllPersisted() {
	reviews := []createReviewBody{
		{Author: "Alice", Rating: 5, Message: "Excellent"},
		{Author: "Bob", Rating: 3, Message: "Average"},
		{Author: "Carol", Rating: 1, Message: "Poor"},
	}

	for _, r := range reviews {
		s.mustPostReview(r)
	}

	s.Equal(int64(len(reviews)), s.dbReviewCount())
}

func (s *ReviewSuite) TestPostReview_ValidationErrors() {
	longMessage := string(make([]byte, 101))

	cases := []struct {
		name       string
		body       createReviewBody
		wantStatus int
	}{
		{
			name:       "Rating above maximum (6)",
			body:       createReviewBody{Author: "Alice", Rating: 6, Message: "msg"},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "Rating below minimum (0)",
			body:       createReviewBody{Author: "Alice", Rating: 0, Message: "msg"},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "Rating is negative",
			body:       createReviewBody{Author: "Alice", Rating: -1, Message: "msg"},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "Author exceeds maxLength (11 chars)",
			body:       createReviewBody{Author: "TooLongName", Rating: 5, Message: "msg"},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "Author exceeds maxLength (20 chars)",
			body:       createReviewBody{Author: "ThisNameIsWayTooLong", Rating: 5, Message: "msg"},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "Message exceeds maxLength (101 chars)",
			body:       createReviewBody{Author: "Alice", Rating: 5, Message: longMessage},
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			testhelper.CleanMongoCollection(s.T(), s.MongoDB.Collection("reviews"))

			w := s.postReview(tc.body)

			s.Equal(tc.wantStatus, w.Code, "unexpected status for case %q", tc.name)
			s.NotEqual(http.StatusCreated, w.Code, "invalid payload must not create a resource")

			s.Equal(int64(0), s.dbReviewCount(), "DB must stay empty after failed POST for case %q", tc.name)
		})
	}
}

func (s *ReviewSuite) TestPostReview_BoundaryValues() {
	exactlyTenChars := "1234567890"
	exactly100Chars := string(make([]byte, 100))

	cases := []struct {
		name string
		body createReviewBody
	}{
		{
			name: "Rating at minimum (1)",
			body: createReviewBody{Author: "Alice", Rating: 1, Message: "msg"},
		},
		{
			name: "Rating at maximum (5)",
			body: createReviewBody{Author: "Alice", Rating: 5, Message: "msg"},
		},
		{
			name: "Author at maxLength (10 chars)",
			body: createReviewBody{Author: exactlyTenChars, Rating: 3, Message: "msg"},
		},
		{
			name: "Message at maxLength (100 chars)",
			body: createReviewBody{Author: "Alice", Rating: 3, Message: exactly100Chars},
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			testhelper.CleanMongoCollection(s.T(), s.MongoDB.Collection("reviews"))

			w := s.postReview(tc.body)

			s.Equal(http.StatusCreated, w.Code, "boundary value should be accepted for case %q", tc.name)
			s.Equal(int64(1), s.dbReviewCount(), "one document should be persisted for case %q", tc.name)
		})
	}
}

func (s *ReviewSuite) TestGetReviews_EmptyDatabase_ReturnsEmptyList() {
	w := s.getReviews()

	s.Equal(http.StatusOK, w.Code)

	resp := s.decodeGetReviewsResponse(w)
	s.NotNil(resp.Data, "data field should be present even when empty")
	s.Empty(resp.Data)
}

func (s *ReviewSuite) TestGetReviews_ReturnsCorrectFields() {
	now := time.Now().UTC()
	coll := s.MongoDB.Collection("reviews")

	body := validReview()

	_, err := coll.InsertOne(
		context.Background(),
		bson.M{"author": body.Author, "rating": body.Rating, "message": body.Message, "createdAt": now},
	)
	s.Require().NoError(err)

	w := s.getReviews()
	s.Equal(http.StatusOK, w.Code)

	resp := s.decodeGetReviewsResponse(w)
	s.Require().Len(resp.Data, 1)

	r := resp.Data[0]
	s.Equal(body.Author, r.Author)
	s.Equal(body.Rating, r.Rating)
	s.Equal(body.Message, r.Message)
	s.NotEmpty(r.ID, "id should be populated (readOnly field)")
	s.False(r.CreatedAt.IsZero(), "createdAt should be populated (readOnly field)")
}

func (s *ReviewSuite) TestGetReviews_ReturnsMostRecentFirst() {
	now := time.Now().UTC()
	coll := s.MongoDB.Collection("reviews")

	_, err := coll.InsertMany(context.Background(), []any{
		bson.M{"author": "Oldest", "rating": 3, "message": "", "createdAt": now.Add(-2 * time.Hour)},
		bson.M{"author": "Middle", "rating": 3, "message": "", "createdAt": now.Add(-1 * time.Hour)},
		bson.M{"author": "Newest", "rating": 3, "message": "", "createdAt": now},
	})
	s.Require().NoError(err)

	w := s.getReviews()
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetReviewsResponse(w)
	s.Require().Len(resp.Data, 3)
	s.Equal("Newest", resp.Data[0].Author)
	s.Equal("Middle", resp.Data[1].Author)
	s.Equal("Oldest", resp.Data[2].Author)
}

func (s *ReviewSuite) TestPostThenGetReviews_ReviewAppearsInList() {
	body := validReview()
	s.mustPostReview(body)

	w := s.getReviews()
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetReviewsResponse(w)
	s.Require().NotEmpty(resp.Data)

	r := resp.Data[0]
	s.Equal(body.Author, r.Author)
	s.Equal(body.Rating, r.Rating)
	s.Equal(body.Message, r.Message)
	s.NotEmpty(r.ID)
	s.False(r.CreatedAt.IsZero())
}

func (s *ReviewSuite) TestPostThenGetReviews_MultipleReviews_AllReturnedInOrder() {
	authors := []string{"First", "Second", "Third"}
	for i, author := range authors {
		s.mustPostReview(createReviewBody{Author: author, Rating: i + 1, Message: "msg"})
	}

	w := s.getReviews()
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetReviewsResponse(w)
	s.Require().Len(resp.Data, len(authors))

	returnedAuthors := make([]string, len(resp.Data))
	for i, r := range resp.Data {
		returnedAuthors[i] = r.Author
	}
	for _, a := range authors {
		s.Contains(returnedAuthors, a)
	}
}
