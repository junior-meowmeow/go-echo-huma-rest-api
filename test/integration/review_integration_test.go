package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/test/helper/adaptor/database"
)

type ReviewSuite struct {
	IntegrationTestSuite
}

func TestReviewSuite(t *testing.T) {
	suite.Run(t, new(ReviewSuite))
}

func (s *ReviewSuite) SetupSuite() {
	s.SetupMongo()
}

func (s *ReviewSuite) SetupTest() {
	s.Database.CleanReviews(s.T())
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

func (s *ReviewSuite) TestPostReview_CreatesReview_PersistsToDatabase() {
	body := validReview()
	s.mustPostReview(body)

	s.Equal(int64(1), s.Database.CountReviews(s.T()))

	records := s.Database.GetAllReviews(s.T())
	s.Require().Len(records, 1)
	doc := records[0]

	s.Equal(body.Author, doc.Author)
	s.Equal(body.Rating, doc.Rating)
	s.Equal(body.Message, doc.Message)
	s.NotEmpty(doc.ID)
	s.NotNil(doc.CreatedAt)
	s.NotNil(doc.UpdatedAt)
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

	s.Equal(int64(len(reviews)), s.Database.CountReviews(s.T()))
}

func (s *ReviewSuite) TestPostReview_ValidationErrors() {
	longMessage := strings.Repeat("a", 101)

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
			s.Database.CleanReviews(s.T())

			w := s.postReview(tc.body)

			s.Equal(tc.wantStatus, w.Code, "unexpected status for case %q", tc.name)
			s.NotEqual(http.StatusCreated, w.Code, "invalid payload must not create a resource")

			s.Equal(int64(0), s.Database.CountReviews(s.T()), "DB must stay empty after failed POST for case %q", tc.name)
		})
	}
}

func (s *ReviewSuite) TestPostReview_BoundaryValues() {
	exactlyTenChars := "1234567890"
	exactly100Chars := strings.Repeat("a", 100)

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
			s.Database.CleanReviews(s.T())

			w := s.postReview(tc.body)

			s.Equal(http.StatusCreated, w.Code, "boundary value should be accepted for case %q", tc.name)
			s.Equal(int64(1), s.Database.CountReviews(s.T()), "one document should be persisted for case %q", tc.name)
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

	body := validReview()

	s.Database.SeedReviews(s.T(), []database.TestReviewRecord{
		{Author: body.Author, Rating: body.Rating, Message: body.Message, CreatedAt: now, UpdatedAt: now},
	})

	w := s.getReviews()
	s.Equal(http.StatusOK, w.Code)

	resp := s.decodeGetReviewsResponse(w)
	s.Require().Len(resp.Data, 1)

	r := resp.Data[0]
	s.Equal(body.Author, r.Author)
	s.Equal(body.Rating, r.Rating)
	s.Equal(body.Message, r.Message)
	s.NotEmpty(r.ID)
	s.False(r.CreatedAt.IsZero())
}

func (s *ReviewSuite) TestGetReviews_ReturnsMostRecentFirst() {
	now := time.Now().UTC()

	s.Database.SeedReviews(s.T(), []database.TestReviewRecord{
		{Author: "Oldest", Rating: 3, Message: "", CreatedAt: now.Add(-2 * time.Hour)},
		{Author: "Middle", Rating: 3, Message: "", CreatedAt: now.Add(-1 * time.Hour)},
		{Author: "Newest", Rating: 3, Message: "", CreatedAt: now},
	})

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
