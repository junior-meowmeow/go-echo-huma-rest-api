package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookSuite struct {
	IntegrationTestSuite
}

func TestBookSuite(t *testing.T) {
	suite.Run(t, new(BookSuite))
}

func (s *BookSuite) SetupTest() {
	cleanCollection(s.T(), s.MongoDB.Collection("books"))
}

type bookMetadataBody struct {
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
	Genre  string `json:"genre,omitempty"`
}

type createBookBody struct {
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Metadata         bookMetadataBody `json:"metadata"`
	CoverImageFileID string           `json:"coverImageFileId,omitempty"`
}

type createBookResponse struct {
	ID string `json:"id"`
}

type bookSchema struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Metadata         bookMetadataBody `json:"metadata"`
	CoverImageFileID string           `json:"coverImageFileId,omitempty"`
	CreatedAt        time.Time        `json:"createdAt"`
}

type getBooksResponse struct {
	Data []bookSchema `json:"data"`
}

func (s *BookSuite) postBook(body createBookBody) *httptest.ResponseRecorder {
	s.T().Helper()

	payload, err := json.Marshal(body)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/books", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *BookSuite) getBooks(queryParams string) *httptest.ResponseRecorder {
	s.T().Helper()

	url := "/v1/books"
	if queryParams != "" {
		url += "?" + queryParams
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *BookSuite) getBookByID(id string) *httptest.ResponseRecorder {
	s.T().Helper()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/books/"+id, nil)
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *BookSuite) mustPostBook(body createBookBody) string {
	s.T().Helper()

	w := s.postBook(body)
	if w.Code != http.StatusCreated {
		s.T().Logf("mustPostBook failed - status: %d, body: %s", w.Code, w.Body.String())
	}
	s.Require().Equal(http.StatusCreated, w.Code)

	var resp createBookResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err, "failed to decode CreateBook response body")
	s.Require().NotEmpty(resp.ID, "CreateBook response must include an ID")

	return resp.ID
}

func (s *BookSuite) decodeGetBooksResponse(w *httptest.ResponseRecorder) getBooksResponse {
	s.T().Helper()

	var resp getBooksResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err, "failed to decode GET /v1/books response")
	return resp
}

func (s *BookSuite) decodeGetBookByIDResponse(w *httptest.ResponseRecorder) bookSchema {
	s.T().Helper()

	var resp bookSchema
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err, "failed to decode GET /v1/books/{id} response")
	return resp
}

func (s *BookSuite) dbBookCount() int64 {
	s.T().Helper()

	count, err := s.MongoDB.Collection("books").CountDocuments(context.Background(), bson.D{})
	s.Require().NoError(err)
	return count
}

func validBook() createBookBody {
	return createBookBody{
		Name:        "Test Book",
		Description: "A book for integration test.",
		Metadata: bookMetadataBody{
			Author: "Junior MeowMeow",
			ISBN:   "978-0132350884",
			Genre:  "Software Engineering",
		},
	}
}

func (s *BookSuite) TestPostBook_ReturnsHTTP201() {
	w := s.postBook(validBook())

	s.Equal(http.StatusCreated, w.Code)
}

func (s *BookSuite) TestPostBook_ReturnsCreatedID() {
	w := s.postBook(validBook())
	s.Require().Equal(http.StatusCreated, w.Code)

	var resp createBookResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)

	s.NotEmpty(resp.ID, "response body must include the created book's ID")
	err = uuid.Validate(resp.ID)
	s.Require().NoError(err, "ID must be a valid UUID")
}

func (s *BookSuite) TestPostBook_PersistsToMongoDB() {
	body := validBook()
	s.mustPostBook(body)

	s.Equal(int64(1), s.dbBookCount())

	var doc bson.M
	err := s.MongoDB.Collection("books").FindOne(context.Background(), bson.D{}).Decode(&doc)
	s.Require().NoError(err)

	s.Equal(body.Name, doc["name"])
	s.Equal(body.Description, doc["description"])
	s.NotEmpty(doc["_id"])
	s.NotEmpty(doc["createdAt"])

	meta, ok := doc["metadata"].(bson.D)
	s.Require().True(ok, "metadata should be stored as a nested document")

	metaMap := make(map[string]any)
	for _, elem := range meta {
		metaMap[elem.Key] = elem.Value
	}

	s.Equal(body.Metadata.Author, metaMap["author"])
	s.Equal(body.Metadata.ISBN, metaMap["isbn"])
	s.Equal(body.Metadata.Genre, metaMap["genre"])
}

func (s *BookSuite) TestPostBook_WithOptionalCoverImage() {
	body := validBook()
	body.CoverImageFileID = "507f1f77bcf86cd799439011"

	id := s.mustPostBook(body)
	s.NotEmpty(id)

	var doc bson.M
	err := s.MongoDB.Collection("books").FindOne(context.Background(), bson.D{}).Decode(&doc)
	s.Require().NoError(err)
	s.Equal(body.CoverImageFileID, doc["coverImageFileId"])
}

func (s *BookSuite) TestPostBook_WithoutOptionalFields() {
	body := createBookBody{
		Name:        "Minimal Book",
		Description: "No genre, no cover.",
		Metadata: bookMetadataBody{
			Author: "Alice",
			ISBN:   "000-0000000000",
			// Genre omitted
		},
		// CoverImageFileID omitted
	}

	w := s.postBook(body)
	s.Equal(http.StatusCreated, w.Code)
	s.Equal(int64(1), s.dbBookCount())
}

func (s *BookSuite) TestPostBook_MultipleBooks_AllPersisted() {
	books := []createBookBody{
		{Name: "Book A", Description: "Desc A", Metadata: bookMetadataBody{Author: "Author A", ISBN: "111"}},
		{Name: "Book B", Description: "Desc B", Metadata: bookMetadataBody{Author: "Author B", ISBN: "222"}},
		{Name: "Book C", Description: "Desc C", Metadata: bookMetadataBody{Author: "Author C", ISBN: "333"}},
	}

	ids := make(map[string]bool)
	for _, b := range books {
		id := s.mustPostBook(b)
		ids[id] = true
	}

	s.Equal(int64(len(books)), s.dbBookCount())
	s.Len(ids, len(books), "every created book must have a unique ID")
}

func (s *BookSuite) TestPostBook_ValidationErrors() {
	longName := string(make([]byte, 101))
	longDescription := string(make([]byte, 501))

	cases := []struct {
		name       string
		body       createBookBody
		wantStatus int
	}{
		{
			name:       "Name exceeds maxLength (101 chars)",
			body:       createBookBody{Name: longName, Description: "ok", Metadata: bookMetadataBody{Author: "A", ISBN: "1"}},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "Description exceeds maxLength (501 chars)",
			body:       createBookBody{Name: "ok", Description: longDescription, Metadata: bookMetadataBody{Author: "A", ISBN: "1"}},
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			cleanCollection(s.T(), s.MongoDB.Collection("books"))

			w := s.postBook(tc.body)

			s.Equal(tc.wantStatus, w.Code, "case: %q", tc.name)
			s.NotEqual(http.StatusCreated, w.Code, "invalid payload must not create a resource")
			s.Equal(int64(0), s.dbBookCount(), "DB must stay empty after failed POST for case %q", tc.name)
		})
	}
}

func (s *BookSuite) TestPostBook_BoundaryValues() {
	exactly100Chars := string(make([]byte, 100))
	exactly500Chars := string(make([]byte, 500))

	cases := []struct {
		name string
		body createBookBody
	}{
		{
			name: "Name at maxLength (100 chars)",
			body: createBookBody{Name: exactly100Chars, Description: "ok", Metadata: bookMetadataBody{Author: "A", ISBN: "1"}},
		},
		{
			name: "Description at maxLength (500 chars)",
			body: createBookBody{Name: "ok", Description: exactly500Chars, Metadata: bookMetadataBody{Author: "A", ISBN: "1"}},
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			cleanCollection(s.T(), s.MongoDB.Collection("books"))

			w := s.postBook(tc.body)

			s.Equal(http.StatusCreated, w.Code, "boundary value should be accepted for case %q", tc.name)
			s.Equal(int64(1), s.dbBookCount(), "one document should be persisted for case %q", tc.name)
		})
	}
}

func (s *BookSuite) TestGetBooks_EmptyDatabase_ReturnsEmptyList() {
	w := s.getBooks("all=true")

	s.Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBooksResponse(w)
	s.NotNil(resp.Data)
	s.Empty(resp.Data)
}

func (s *BookSuite) TestGetBooks_ReturnsCorrectFields() {
	body := validBook()
	s.mustPostBook(body)

	w := s.getBooks("all=true")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBooksResponse(w)
	s.Require().Len(resp.Data, 1)

	b := resp.Data[0]
	s.Equal(body.Name, b.Name)
	s.Equal(body.Description, b.Description)
	s.Equal(body.Metadata.Author, b.Metadata.Author)
	s.Equal(body.Metadata.ISBN, b.Metadata.ISBN)
	s.Equal(body.Metadata.Genre, b.Metadata.Genre)
	s.NotEmpty(b.ID, "id should be populated (readOnly field)")
	s.False(b.CreatedAt.IsZero(), "createdAt should be populated (readOnly field)")
}

func (s *BookSuite) TestGetBooks_ReturnsMostRecentFirst() {
	now := time.Now().UTC()
	coll := s.MongoDB.Collection("books")

	_, err := coll.InsertMany(context.Background(), []any{
		bson.M{"_id": uuid.New(), "name": "Oldest", "description": "", "metadata": bson.M{}, "createdAt": now.Add(-2 * time.Hour)},
		bson.M{"_id": uuid.New(), "name": "Middle", "description": "", "metadata": bson.M{}, "createdAt": now.Add(-1 * time.Hour)},
		bson.M{"_id": uuid.New(), "name": "Newest", "description": "", "metadata": bson.M{}, "createdAt": now},
	})
	s.Require().NoError(err)

	w := s.getBooks("all=true")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBooksResponse(w)
	s.Require().Len(resp.Data, 3)
	s.Equal("Newest", resp.Data[0].Name)
	s.Equal("Middle", resp.Data[1].Name)
	s.Equal("Oldest", resp.Data[2].Name)
}

func (s *BookSuite) TestGetBooks_Pagination_DefaultPage() {
	docs := make([]any, 25)
	for i := range docs {
		docs[i] = bson.M{
			"_id":       uuid.New(),
			"name":      fmt.Sprintf("Book %02d", i),
			"metadata":  bson.M{},
			"createdAt": time.Now().UTC(),
		}
	}
	_, err := s.MongoDB.Collection("books").InsertMany(context.Background(), docs)
	s.Require().NoError(err)

	w := s.getBooks("")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBooksResponse(w)
	s.Len(resp.Data, 20, "default page size is 20")
}

func (s *BookSuite) TestGetBooks_Pagination_SecondPage() {
	docs := make([]any, 25)
	for i := range docs {
		docs[i] = bson.M{
			"_id":       uuid.New(),
			"name":      fmt.Sprintf("Book %02d", i),
			"metadata":  bson.M{},
			"createdAt": time.Now().UTC(),
		}
	}
	_, err := s.MongoDB.Collection("books").InsertMany(context.Background(), docs)
	s.Require().NoError(err)

	w := s.getBooks("pageNumber=2&pageSize=20")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBooksResponse(w)
	s.Len(resp.Data, 5, "page 2 should contain the remaining 5 books")
}

func (s *BookSuite) TestGetBooks_Pagination_CustomPageSize() {
	docs := make([]any, 10)
	for i := range docs {
		docs[i] = bson.M{"_id": uuid.New(), "name": fmt.Sprintf("Book %02d", i), "metadata": bson.M{}, "createdAt": time.Now().UTC()}
	}
	_, err := s.MongoDB.Collection("books").InsertMany(context.Background(), docs)
	s.Require().NoError(err)

	w := s.getBooks("pageSize=3&pageNumber=1")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBooksResponse(w)
	s.Len(resp.Data, 3)
}

func (s *BookSuite) TestGetBooks_Pagination_GetAll_IgnoresPagination() {
	docs := make([]any, 25)
	for i := range docs {
		docs[i] = bson.M{"_id": uuid.New(), "name": fmt.Sprintf("Book %02d", i), "metadata": bson.M{}, "createdAt": time.Now().UTC()}
	}
	_, err := s.MongoDB.Collection("books").InsertMany(context.Background(), docs)
	s.Require().NoError(err)

	w := s.getBooks("all=true&pageSize=3")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBooksResponse(w)
	s.Len(resp.Data, 25, "all=true must return all books regardless of pageSize")
}

func (s *BookSuite) TestGetBooks_Pagination_ValidationErrors() {
	cases := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{
			name:       "pageNumber below minimum (0)",
			query:      "pageNumber=0",
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "pageSize below minimum (0)",
			query:      "pageSize=0",
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "pageSize above maximum (101)",
			query:      "pageSize=101",
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			w := s.getBooks(tc.query)
			s.Equal(tc.wantStatus, w.Code, "case: %q", tc.name)
		})
	}
}

func (s *BookSuite) TestGetBookByID_ReturnsCorrectBook() {
	body := validBook()
	id := s.mustPostBook(body)

	w := s.getBookByID(id)
	s.Require().Equal(http.StatusOK, w.Code)

	book := s.decodeGetBookByIDResponse(w)
	s.Equal(id, book.ID)
	s.Equal(body.Name, book.Name)
	s.Equal(body.Description, book.Description)
	s.Equal(body.Metadata.Author, book.Metadata.Author)
	s.Equal(body.Metadata.ISBN, book.Metadata.ISBN)
	s.Equal(body.Metadata.Genre, book.Metadata.Genre)
	s.False(book.CreatedAt.IsZero())
}

func (s *BookSuite) TestGetBookByID_NotFound_Returns404() {
	nonExistentID := uuid.NewString()

	w := s.getBookByID(nonExistentID)

	s.Equal(http.StatusNotFound, w.Code)
}

func (s *BookSuite) TestGetBookByID_InvalidIDFormat_Returns422() {
	cases := []struct {
		name string
		id   string
	}{
		{name: "Too short", id: "abc123"},
		{name: "Non-hex characters", id: "zzzzzzzzzzzzzzzzzzzzzzzz"},
		{name: "Too long", id: "507f1f77bcf86cd7994390111"},
		{name: "Empty-ish (spaces)", id: "   "},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			w := s.getBookByID(tc.id)
			s.Equal(http.StatusUnprocessableEntity, w.Code, "case: %q", tc.name)
		})
	}
}

func (s *BookSuite) TestPostThenGetBooks_BookAppearsInList() {
	body := validBook()
	id := s.mustPostBook(body)

	w := s.getBooks("all=true")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBooksResponse(w)
	s.Require().NotEmpty(resp.Data)

	b := resp.Data[0]
	s.Equal(id, b.ID)
	s.Equal(body.Name, b.Name)
	s.Equal(body.Metadata.Author, b.Metadata.Author)
}

func (s *BookSuite) TestPostThenGetBookByID_RoundTrip() {
	body := validBook()
	id := s.mustPostBook(body)

	w := s.getBookByID(id)
	s.Require().Equal(http.StatusOK, w.Code)

	book := s.decodeGetBookByIDResponse(w)
	s.Equal(id, book.ID)
	s.Equal(body.Name, book.Name)
	s.Equal(body.Description, book.Description)
	s.Equal(body.Metadata.Author, book.Metadata.Author)
	s.Equal(body.Metadata.ISBN, book.Metadata.ISBN)
	s.Equal(body.Metadata.Genre, book.Metadata.Genre)
	s.False(book.CreatedAt.IsZero())
}

func (s *BookSuite) TestPostThenGetBooks_MultipleBooks_AllPresent() {
	books := []createBookBody{
		{Name: "Book A", Description: "Desc A", Metadata: bookMetadataBody{Author: "Author A", ISBN: "111"}},
		{Name: "Book B", Description: "Desc B", Metadata: bookMetadataBody{Author: "Author B", ISBN: "222"}},
		{Name: "Book C", Description: "Desc C", Metadata: bookMetadataBody{Author: "Author C", ISBN: "333"}},
	}

	postedIDs := make(map[string]bool)
	for _, b := range books {
		postedIDs[s.mustPostBook(b)] = true
	}

	w := s.getBooks("all=true")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBooksResponse(w)
	s.Require().Len(resp.Data, len(books))

	returnedIDs := make(map[string]bool)
	for _, b := range resp.Data {
		returnedIDs[b.ID] = true
	}
	s.Equal(postedIDs, returnedIDs, "all posted IDs should appear in the list")
}
