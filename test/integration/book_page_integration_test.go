package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/test/helper/adaptor/database"
)

type BookPageSuite struct {
	IntegrationTestSuite

	// bookID is a valid book created once per test, used as the parent for
	// all book pages. Re-created in SetupTest so every test is isolated.
	bookID string
}

func (s *BookPageSuite) SetupTest() {
	s.Database.CleanBookPages(s.T())
	s.Database.CleanBooks(s.T())

	s.bookID = uuid.NewString()
	parentBook := database.TestBookRecord{
		ID:          s.bookID,
		Name:        "Test Book",
		Description: "Parent book for book page tests.",
		Author:      "Junior MeowMeow",
		ISBN:        "978-0132350884",
		Genre:       "Test",
		CreatedAt:   time.Now().UTC(),
	}
	s.Database.SeedBooks(s.T(), []database.TestBookRecord{parentBook})
}

type MongoBookPageSuite struct{ BookPageSuite }

func (s *MongoBookPageSuite) SetupSuite() {
	s.SetupMongo()
}

func TestMongoBookPageSuite(t *testing.T) {
	suite.Run(t, new(MongoBookPageSuite))
}

type PostgresBookPageSuite struct{ BookPageSuite }

func (s *PostgresBookPageSuite) SetupSuite() {
	s.SetupPostgres()
}

func TestPostgresBookPageSuite(t *testing.T) {
	suite.Run(t, new(PostgresBookPageSuite))
}

type bookPageMetadataBody struct {
	IsBookmarked bool   `json:"isBookmarked"`
	Highlight    string `json:"highlight,omitempty"`
}

type createBookPageBody struct {
	BookID              string               `json:"bookId"`
	PageNumber          int64                `json:"pageNumber"`
	Content             string               `json:"content"`
	Metadata            bookPageMetadataBody `json:"metadata"`
	AttachedImageFileID string               `json:"attachedImageFileId,omitempty"`
}

type createBookPageResponse struct {
	ID string `json:"id"`
}

type bookPageSchema struct {
	ID                  string               `json:"id"`
	BookID              string               `json:"bookId"`
	PageNumber          int64                `json:"pageNumber"`
	Content             string               `json:"content"`
	Metadata            bookPageMetadataBody `json:"metadata"`
	AttachedImageFileID string               `json:"attachedImageFileId,omitempty"`
}

type getBookPagesResponse struct {
	Data []bookPageSchema `json:"data"`
}

func (s *BookPageSuite) postBookPage(body createBookPageBody) *httptest.ResponseRecorder {
	s.T().Helper()

	payload, err := json.Marshal(body)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/book_pages", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *BookPageSuite) mustPostBookPage(body createBookPageBody) string {
	s.T().Helper()

	w := s.postBookPage(body)
	if w.Code != http.StatusCreated {
		s.T().Logf("mustPostBookPage failed - status: %d, body: %s", w.Code, w.Body.String())
	}
	s.Require().Equal(http.StatusCreated, w.Code)

	var resp createBookPageResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err, "failed to decode CreateBookPage response body")
	s.Require().NotEmpty(resp.ID, "CreateBookPage response must include an ID")

	return resp.ID
}

func (s *BookPageSuite) getBookPages(bookID string, extraQuery string) *httptest.ResponseRecorder {
	s.T().Helper()

	url := fmt.Sprintf("/v1/book_pages?bookId=%s", bookID)
	if extraQuery != "" {
		url += "&" + extraQuery
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *BookPageSuite) getBookPagesByRange(bookID string, startPage, endPage int64) *httptest.ResponseRecorder {
	s.T().Helper()

	url := fmt.Sprintf("/v1/book_pages/range?bookId=%s&startPage=%d&endPage=%d", bookID, startPage, endPage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *BookPageSuite) getBookPagesByOffset(bookID string, centerPage, offset int64) *httptest.ResponseRecorder {
	s.T().Helper()

	url := fmt.Sprintf("/v1/book_pages/offset?bookId=%s&centerPage=%d&offset=%d", bookID, centerPage, offset)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *BookPageSuite) getBookPageByID(id string) *httptest.ResponseRecorder {
	s.T().Helper()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/book_pages/"+id, nil)
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *BookPageSuite) decodeGetBookPagesResponse(w *httptest.ResponseRecorder) getBookPagesResponse {
	s.T().Helper()

	var resp getBookPagesResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err, "failed to decode book pages list response")
	return resp
}

func (s *BookPageSuite) decodeGetBookPageByIDResponse(w *httptest.ResponseRecorder) bookPageSchema {
	s.T().Helper()

	var resp bookPageSchema
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err, "failed to decode GET /v1/book_pages/{id} response")
	return resp
}

func (s *BookPageSuite) validBookPage(bookPageNumber int64) createBookPageBody {
	return createBookPageBody{
		BookID:     s.bookID,
		PageNumber: bookPageNumber,
		Content:    fmt.Sprintf("Content of page %d.", bookPageNumber),
		Metadata: bookPageMetadataBody{
			IsBookmarked: false,
		},
	}
}

func (s *BookPageSuite) seedPages(count int) {
	s.T().Helper()

	pages := make([]database.TestBookPageRecord, count)
	for i := range pages {
		pages[i] = database.TestBookPageRecord{
			ID:           uuid.NewString(),
			BookID:       s.bookID,
			PageNumber:   int64(i + 1),
			Content:      fmt.Sprintf("Page %d content", i+1),
			IsBookmarked: false,
			Highlight:    "",
		}
	}

	s.Database.SeedBookPages(s.T(), pages)
}

func (s *BookPageSuite) TestPostBookPage_ReturnsHTTP201() {
	w := s.postBookPage(s.validBookPage(1))

	s.Equal(http.StatusCreated, w.Code)
}

func (s *BookPageSuite) TestPostBookPage_ReturnsCreatedID() {
	w := s.postBookPage(s.validBookPage(1))
	s.Require().Equal(http.StatusCreated, w.Code)

	var resp createBookPageResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().NoError(err)

	s.NotEmpty(resp.ID, "response body must include the created book page's ID")
	err = uuid.Validate(resp.ID)
	s.Require().NoError(err, "ID must be a valid UUID")
}

func (s *BookPageSuite) TestPostBookPage_PersistsToDatabase() {
	body := s.validBookPage(1)
	id := s.mustPostBookPage(body)

	s.Equal(int64(1), s.Database.CountBookPages(s.T()))

	doc := s.Database.GetBookPageByID(s.T(), id)

	s.Equal(body.PageNumber, doc.PageNumber)
	s.Equal(body.Content, doc.Content)
	s.Equal(s.bookID, doc.BookID)
	s.Equal(body.Metadata.IsBookmarked, doc.IsBookmarked)
}

func (s *BookPageSuite) TestPostBookPage_WithOptionalAttachedImage() {
	body := s.validBookPage(1)
	body.AttachedImageFileID = "507f1f77bcf86cd799439011"

	id := s.mustPostBookPage(body)
	s.NotEmpty(id)

	doc := s.Database.GetBookPageByID(s.T(), id)
	s.Equal(body.AttachedImageFileID, doc.AttachedImageFileID)
}

func (s *BookPageSuite) TestPostBookPage_WithBookmark() {
	body := s.validBookPage(1)
	body.Metadata.IsBookmarked = true
	body.Metadata.Highlight = "Important passage"

	id := s.mustPostBookPage(body)

	doc := s.Database.GetBookPageByID(s.T(), id)
	s.True(doc.IsBookmarked)
	s.Equal("Important passage", doc.Highlight)
}

func (s *BookPageSuite) TestPostBookPage_MultiplePages_AllPersisted() {
	ids := make(map[string]bool)
	for pageNum := int64(1); pageNum <= 5; pageNum++ {
		id := s.mustPostBookPage(s.validBookPage(pageNum))
		ids[id] = true
	}

	s.Equal(int64(5), s.Database.CountBookPages(s.T()))
	s.Len(ids, 5, "every created page must have a unique ID")
}

func (s *BookPageSuite) TestPostBookPage_InvalidBookID_Returns422() {
	cases := []struct {
		name   string
		bookID string
	}{
		{name: "Too short", bookID: "abc123"},
		{name: "Non-hex", bookID: "zzzzzzzzzzzzzzzzzzzzzzzz"},
		{name: "Empty", bookID: ""},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.Database.CleanBookPages(s.T())

			body := s.validBookPage(1)
			body.BookID = tc.bookID
			w := s.postBookPage(body)

			s.Equal(http.StatusUnprocessableEntity, w.Code, "case: %q", tc.name)
			s.Equal(int64(0), s.Database.CountBookPages(s.T()))
		})
	}
}

func (s *BookPageSuite) TestGetBookPages_EmptyBook_ReturnsEmptyList() {
	w := s.getBookPages(s.bookID, "all=true")

	s.Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.NotNil(resp.Data)
	s.Empty(resp.Data)
}

func (s *BookPageSuite) TestGetBookPages_ReturnsCorrectFields() {
	body := s.validBookPage(1)
	body.Metadata.IsBookmarked = true
	body.Metadata.Highlight = "Key insight"
	s.mustPostBookPage(body)

	w := s.getBookPages(s.bookID, "all=true")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.Require().Len(resp.Data, 1)

	p := resp.Data[0]
	s.Equal(s.bookID, p.BookID)
	s.Equal(body.PageNumber, p.PageNumber)
	s.Equal(body.Content, p.Content)
	s.Equal(body.Metadata.IsBookmarked, p.Metadata.IsBookmarked)
	s.Equal(body.Metadata.Highlight, p.Metadata.Highlight)
	s.NotEmpty(p.ID, "id should be populated (readOnly field)")
}

func (s *BookPageSuite) TestGetBookPages_ReturnsPageNumberAscending() {
	// Insert in reverse order.
	for _, pageNum := range []int64{3, 1, 2} {
		s.mustPostBookPage(s.validBookPage(pageNum))
	}

	w := s.getBookPages(s.bookID, "all=true")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.Require().Len(resp.Data, 3)
	s.Equal(int64(1), resp.Data[0].PageNumber)
	s.Equal(int64(2), resp.Data[1].PageNumber)
	s.Equal(int64(3), resp.Data[2].PageNumber)
}

func (s *BookPageSuite) TestGetBookPages_IsolatedByBookID_OnlyReturnsOwnPages() {
	// Insert a second book.
	otherBookID := uuid.NewString()
	otherBook := database.TestBookRecord{
		ID:        otherBookID,
		Name:      "Other Book",
		CreatedAt: time.Now().UTC(),
	}
	s.Database.SeedBooks(s.T(), []database.TestBookRecord{otherBook})

	s.mustPostBookPage(s.validBookPage(1)) // belongs to s.bookID

	otherPage := s.validBookPage(1)
	otherPage.BookID = otherBookID
	s.mustPostBookPage(otherPage) // belongs to otherBookID

	w := s.getBookPages(s.bookID, "all=true")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.Require().Len(resp.Data, 1, "should only return pages belonging to the requested book")
	s.Equal(s.bookID, resp.Data[0].BookID)
}

func (s *BookPageSuite) TestGetBookPages_Pagination_DefaultPage() {
	s.seedPages(60) // default pageSize is 50

	w := s.getBookPages(s.bookID, "")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.Len(resp.Data, 50, "default page size should be 50")
}

func (s *BookPageSuite) TestGetBookPages_Pagination_SecondPage() {
	s.seedPages(60)

	w := s.getBookPages(s.bookID, "pageNumber=2&pageSize=50")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.Len(resp.Data, 10, "second page should contain the remaining 10 pages")
}

func (s *BookPageSuite) TestGetBookPages_Pagination_GetAll_IgnoresPagination() {
	s.seedPages(60)

	w := s.getBookPages(s.bookID, "all=true&pageSize=10")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.Len(resp.Data, 60, "all=true must return all pages regardless of pageSize")
}

func (s *BookPageSuite) TestGetBookPages_Pagination_ValidationErrors() {
	cases := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{
			name:       "pageSize above maximum (501)",
			query:      "pageSize=501",
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "pageSize below minimum (0)",
			query:      "pageSize=0",
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "pageNumber below minimum (0)",
			query:      "pageNumber=0",
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			w := s.getBookPages(s.bookID, tc.query)
			s.Equal(tc.wantStatus, w.Code, "case: %q", tc.name)
		})
	}
}

func (s *BookPageSuite) TestGetBookPages_MissingBookID_Returns422() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/book_pages?all=true", nil)
	s.Router.ServeHTTP(w, req)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
}

func (s *BookPageSuite) TestGetBookPages_InvalidBookID_Returns422() {
	w := s.getBookPages("not-a-valid-objectid", "all=true")
	s.Equal(http.StatusUnprocessableEntity, w.Code)
}

func (s *BookPageSuite) TestGetBookPagesByRange_ReturnsCorrectPages() {
	s.seedPages(10)

	w := s.getBookPagesByRange(s.bookID, 3, 6)
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.Require().Len(resp.Data, 4, "range [3,6] is inclusive on both ends")
	s.Equal(int64(3), resp.Data[0].PageNumber)
	s.Equal(int64(4), resp.Data[1].PageNumber)
	s.Equal(int64(5), resp.Data[2].PageNumber)
	s.Equal(int64(6), resp.Data[3].PageNumber)
}

func (s *BookPageSuite) TestGetBookPagesByRange_SinglePage() {
	s.seedPages(5)

	w := s.getBookPagesByRange(s.bookID, 3, 3)
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.Require().Len(resp.Data, 1)
	s.Equal(int64(3), resp.Data[0].PageNumber)
}

func (s *BookPageSuite) TestGetBookPagesByRange_RangeExceedsAvailable_ReturnsWhatExists() {
	s.seedPages(5)

	w := s.getBookPagesByRange(s.bookID, 4, 10)
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	// Should return pages 4 and 5 only; pages 6-10 don't exist.
	s.Len(resp.Data, 2)
	s.Equal(int64(4), resp.Data[0].PageNumber)
	s.Equal(int64(5), resp.Data[1].PageNumber)
}

func (s *BookPageSuite) TestGetBookPagesByRange_EmptyRange_ReturnsEmptyList() {
	s.seedPages(5)

	// Range entirely outside existing pages.
	w := s.getBookPagesByRange(s.bookID, 20, 30)
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.Empty(resp.Data)
}

func (s *BookPageSuite) TestGetBookPagesByRange_MissingBookID_Returns422() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/book_pages/range?startPage=1&endPage=5", nil)
	s.Router.ServeHTTP(w, req)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
}

func (s *BookPageSuite) TestGetBookPagesByRange_MissingRequiredParams_Returns422() {
	cases := []struct {
		name  string
		query string
	}{
		{name: "Missing startPage", query: fmt.Sprintf("bookId=%s&endPage=5", s.bookID)},
		{name: "Missing endPage", query: fmt.Sprintf("bookId=%s&startPage=1", s.bookID)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/v1/book_pages/range?"+tc.query, nil)
			s.Router.ServeHTTP(w, req)
			s.Equal(http.StatusUnprocessableEntity, w.Code, "case: %q", tc.name)
		})
	}
}

func (s *BookPageSuite) TestGetBookPagesByOffset_ReturnsCorrectWindow() {
	s.seedPages(10)

	w := s.getBookPagesByOffset(s.bookID, 5, 2)
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.Require().Len(resp.Data, 5)

	pageNums := make([]int64, len(resp.Data))
	for i, p := range resp.Data {
		pageNums[i] = p.PageNumber
	}
	s.Equal([]int64{3, 4, 5, 6, 7}, pageNums)
}

func (s *BookPageSuite) TestGetBookPagesByOffset_AtStart_CapsToAvailable() {
	s.seedPages(10)

	w := s.getBookPagesByOffset(s.bookID, 2, 3)
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)

	pageNums := make([]int64, len(resp.Data))
	for i, p := range resp.Data {
		pageNums[i] = p.PageNumber
	}
	s.Equal([]int64{1, 2, 3, 4, 5}, pageNums)
}

func (s *BookPageSuite) TestGetBookPagesByOffset_AtEnd_CapsToAvailable() {
	s.seedPages(10)

	w := s.getBookPagesByOffset(s.bookID, 9, 3)
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)

	pageNums := make([]int64, len(resp.Data))
	for i, p := range resp.Data {
		pageNums[i] = p.PageNumber
	}
	s.Equal([]int64{6, 7, 8, 9, 10}, pageNums)
}

func (s *BookPageSuite) TestGetBookPagesByOffset_ZeroOffset_ReturnsOnlyCenterPage() {
	s.seedPages(5)

	w := s.getBookPagesByOffset(s.bookID, 3, 0)
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.Require().Len(resp.Data, 1)
	s.Equal(int64(3), resp.Data[0].PageNumber)
}

func (s *BookPageSuite) TestGetBookPagesByOffset_MissingBookID_Returns422() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/book_pages/offset?centerPage=5&offset=2", nil)
	s.Router.ServeHTTP(w, req)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
}

func (s *BookPageSuite) TestGetBookPagesByOffset_MissingRequiredParams_Returns422() {
	cases := []struct {
		name  string
		query string
	}{
		{name: "Missing centerPage", query: fmt.Sprintf("bookId=%s&offset=2", s.bookID)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/v1/book_pages/offset?"+tc.query, nil)
			s.Router.ServeHTTP(w, req)
			s.Equal(http.StatusUnprocessableEntity, w.Code, "case: %q", tc.name)
		})
	}
}

func (s *BookPageSuite) TestGetBookPageByID_ReturnsCorrectPage() {
	body := s.validBookPage(7)
	body.Metadata.IsBookmarked = true
	id := s.mustPostBookPage(body)

	w := s.getBookPageByID(id)
	s.Require().Equal(http.StatusOK, w.Code)

	p := s.decodeGetBookPageByIDResponse(w)
	s.Equal(id, p.ID)
	s.Equal(s.bookID, p.BookID)
	s.Equal(body.PageNumber, p.PageNumber)
	s.Equal(body.Content, p.Content)
	s.Equal(body.Metadata.IsBookmarked, p.Metadata.IsBookmarked)
}

func (s *BookPageSuite) TestGetBookPageByID_NotFound_Returns404() {
	nonExistentID := uuid.NewString()

	w := s.getBookPageByID(nonExistentID)

	s.Equal(http.StatusNotFound, w.Code)
}

func (s *BookPageSuite) TestGetBookPageByID_InvalidIDFormat_Returns422() {
	cases := []struct {
		name string
		id   string
	}{
		{name: "Too short", id: "abc123"},
		{name: "Non-hex characters", id: "zzzzzzzzzzzzzzzzzzzzzzzz"},
		{name: "Too long", id: "507f1f77bcf86cd7994390111"},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			w := s.getBookPageByID(tc.id)
			s.Equal(http.StatusUnprocessableEntity, w.Code, "case: %q", tc.name)
		})
	}
}

func (s *BookPageSuite) TestPostThenGetBookPages_PageAppearsInList() {
	body := s.validBookPage(1)
	id := s.mustPostBookPage(body)

	w := s.getBookPages(s.bookID, "all=true")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.Require().Len(resp.Data, 1)
	s.Equal(id, resp.Data[0].ID)
	s.Equal(body.PageNumber, resp.Data[0].PageNumber)
	s.Equal(body.Content, resp.Data[0].Content)
}

func (s *BookPageSuite) TestPostThenGetBookPageByID_RoundTrip() {
	body := s.validBookPage(5)
	id := s.mustPostBookPage(body)

	w := s.getBookPageByID(id)
	s.Require().Equal(http.StatusOK, w.Code)

	p := s.decodeGetBookPageByIDResponse(w)
	s.Equal(id, p.ID)
	s.Equal(s.bookID, p.BookID)
	s.Equal(body.PageNumber, p.PageNumber)
	s.Equal(body.Content, p.Content)
}

func (s *BookPageSuite) TestPostThenGetByRange_PagesAppearInRange() {
	for pageNum := int64(1); pageNum <= 5; pageNum++ {
		s.mustPostBookPage(s.validBookPage(pageNum))
	}

	w := s.getBookPagesByRange(s.bookID, 2, 4)
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.Require().Len(resp.Data, 3)
	s.Equal(int64(2), resp.Data[0].PageNumber)
	s.Equal(int64(3), resp.Data[1].PageNumber)
	s.Equal(int64(4), resp.Data[2].PageNumber)
}

func (s *BookPageSuite) TestPostThenGetByOffset_WindowAroundCenterPage() {
	for pageNum := int64(1); pageNum <= 10; pageNum++ {
		s.mustPostBookPage(s.validBookPage(pageNum))
	}

	w := s.getBookPagesByOffset(s.bookID, 5, 1)
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeGetBookPagesResponse(w)
	s.Require().Len(resp.Data, 3)
	s.Equal(int64(4), resp.Data[0].PageNumber)
	s.Equal(int64(5), resp.Data[1].PageNumber)
	s.Equal(int64(6), resp.Data[2].PageNumber)
}
