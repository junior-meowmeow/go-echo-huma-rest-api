package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type FileSuite struct {
	IntegrationTestSuite
}

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileSuite))
}

func (s *FileSuite) SetupSuite() {
	s.IntegrationTestSuite.SetupSuite()
	s.ensureBucketExists()
}

func (s *FileSuite) SetupTest() {
	cleanCollection(s.T(), s.MongoDB.Collection("filerecords"))
	cleanBucket(s.T(), context.Background(), s.S3Client, "test-bucket")
}

func (s *FileSuite) ensureBucketExists() {
	s.T().Helper()

	ctx := context.Background()
	_, err := s.S3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String("test-bucket"),
	})
	if err != nil && !strings.Contains(err.Error(), "BucketAlreadyExists") && !strings.Contains(err.Error(), "BucketAlreadyOwnedByYou") {
		s.Require().NoError(err, "failed to create test bucket")
	}
}

type uploadFileResponse struct {
	ID string `json:"id"`
}

type fileDownloadInfoResponse struct {
	FileName    string    `json:"fileName"`
	DownloadURL string    `json:"downloadUrl"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

func (s *FileSuite) uploadFile(content []byte, fileName, contentType, objectBaseKey string) *httptest.ResponseRecorder {
	s.T().Helper()

	var body bytes.Buffer
	mw := multipart.NewWriter(&body)

	if objectBaseKey != "" {
		s.Require().NoError(mw.WriteField("objectBaseKey", objectBaseKey))
	}

	h := make(map[string][]string)
	h["Content-Disposition"] = []string{fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName)}
	h["Content-Type"] = []string{contentType}
	part, err := mw.CreatePart(h)
	s.Require().NoError(err, "failed to create form file part")
	_, err = io.Copy(part, bytes.NewReader(content))
	s.Require().NoError(err, "failed to write file content")
	s.Require().NoError(mw.Close())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/files/upload", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *FileSuite) getFileDownloadLink(id string) *httptest.ResponseRecorder {
	s.T().Helper()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/files/download?id="+id, nil)
	s.Router.ServeHTTP(w, req)
	return w
}

func (s *FileSuite) mustUploadFile(content []byte, fileName, contentType, objectBaseKey string) string {
	s.T().Helper()

	w := s.uploadFile(content, fileName, contentType, objectBaseKey)
	if w.Code != http.StatusOK {
		s.T().Logf("mustUploadFile failed – status: %d, body: %s", w.Code, w.Body.String())
	}
	s.Require().Equal(http.StatusOK, w.Code)

	var resp uploadFileResponse
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp), "failed to decode UploadFile response")
	s.Require().NotEmpty(resp.ID, "UploadFile response must include an ID")

	return resp.ID
}

func (s *FileSuite) decodeUploadFileResponse(w *httptest.ResponseRecorder) uploadFileResponse {
	s.T().Helper()

	var resp uploadFileResponse
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp), "failed to decode POST /v1/files/upload response")
	return resp
}

func (s *FileSuite) decodeFileDownloadLinkResponse(w *httptest.ResponseRecorder) fileDownloadInfoResponse {
	s.T().Helper()

	var resp fileDownloadInfoResponse
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp), "failed to decode GET /v1/files/download response")
	return resp
}

func (s *FileSuite) dbFileRecordCount() int64 {
	s.T().Helper()

	count, err := s.MongoDB.Collection("filerecords").CountDocuments(context.Background(), bson.D{})
	s.Require().NoError(err)
	return count
}

func (s *FileSuite) listS3Objects() []string {
	s.T().Helper()

	out, err := s.S3Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String("test-bucket"),
	})
	s.Require().NoError(err)

	keys := make([]string, 0, len(out.Contents))
	for _, obj := range out.Contents {
		keys = append(keys, aws.ToString(obj.Key))
	}
	return keys
}

func (s *FileSuite) TestUploadFile_ReturnsHTTP200() {
	w := s.uploadFile([]byte("hello world"), "test.txt", "text/plain", "")

	s.Equal(http.StatusOK, w.Code)
}

func (s *FileSuite) TestUploadFile_ReturnsCreatedID() {
	w := s.uploadFile([]byte("hello"), "test.txt", "text/plain", "")
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeUploadFileResponse(w)
	s.NotEmpty(resp.ID, "response body must include the created file ID")
	s.Regexp(`^[a-fA-F0-9]{24}$`, resp.ID, "ID must be a valid BSON ObjectID")
}

func (s *FileSuite) TestUploadFile_PersistsFileRecordToMongoDB() {
	s.mustUploadFile([]byte("persist test"), "persist.txt", "text/plain", "")

	s.Equal(int64(1), s.dbFileRecordCount())
}

func (s *FileSuite) TestUploadFile_PersistsObjectToS3() {
	s.mustUploadFile([]byte("s3 content"), "s3test.txt", "text/plain", "")

	keys := s.listS3Objects()
	s.Len(keys, 1, "exactly one object should exist in the bucket after upload")
}

func (s *FileSuite) TestUploadFile_WithObjectBaseKey_UsesKeyAsPrefix() {
	prefix := "uploads/"
	s.mustUploadFile([]byte("prefixed"), "file.txt", "text/plain", prefix)

	keys := s.listS3Objects()
	s.Require().Len(keys, 1)
	s.True(strings.HasPrefix(keys[0], prefix),
		"S3 key %q should start with prefix %q", keys[0], prefix)
}

func (s *FileSuite) TestUploadFile_WithoutObjectBaseKey_KeyHasNoSlash() {
	s.mustUploadFile([]byte("no prefix"), "file.txt", "text/plain", "")

	keys := s.listS3Objects()
	s.Require().Len(keys, 1)
	s.NotContains(keys[0], "/",
		"S3 key %q should not contain a path separator when objectBaseKey is empty", keys[0])
}

func (s *FileSuite) TestUploadFile_PreservesFileExtensionInS3Key() {
	s.mustUploadFile([]byte(`{"x":1}`), "data.json", "application/json", "")

	keys := s.listS3Objects()
	s.Require().Len(keys, 1)
	s.True(strings.HasSuffix(keys[0], ".json"),
		"S3 key %q should preserve the original file extension", keys[0])
}

func (s *FileSuite) TestUploadFile_MultipleFiles_AllPersistedWithUniqueIDs() {
	files := []struct {
		content []byte
		name    string
	}{
		{[]byte("file one"), "one.txt"},
		{[]byte("file two"), "two.txt"},
		{[]byte("file three"), "three.txt"},
	}

	ids := make(map[string]bool)
	for _, f := range files {
		id := s.mustUploadFile(f.content, f.name, "text/plain", "")
		ids[id] = true
	}

	s.Equal(int64(len(files)), s.dbFileRecordCount(), "all file records should be persisted")
	s.Len(ids, len(files), "every uploaded file must have a unique ID")
	s.Len(s.listS3Objects(), len(files), "all files should exist in S3")
}

func (s *FileSuite) TestUploadFile_DifferentContentTypes() {
	cases := []struct {
		name        string
		content     []byte
		fileName    string
		contentType string
	}{
		{"plain text", []byte("hello"), "hello.txt", "text/plain"},
		{"json", []byte(`{"key":"value"}`), "data.json", "application/json"},
		{"binary", []byte{0x89, 0x50, 0x4E, 0x47}, "image.png", "image/png"},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			cleanCollection(s.T(), s.MongoDB.Collection("filerecords"))
			cleanBucket(s.T(), context.Background(), s.S3Client, "test-bucket")

			w := s.uploadFile(tc.content, tc.fileName, tc.contentType, "")
			s.Equal(http.StatusOK, w.Code, "case: %q", tc.name)
			s.Equal(int64(1), s.dbFileRecordCount(), "case: %q", tc.name)
		})
	}
}

func (s *FileSuite) TestUploadFile_NoFilePart_Returns422() {
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	s.Require().NoError(mw.Close())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/files/upload", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	s.Router.ServeHTTP(w, req)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
	s.Equal(int64(0), s.dbFileRecordCount(), "no file record should be created on a failed upload")
}

func (s *FileSuite) TestGetFileDownloadLink_ReturnsHTTP200() {
	id := s.mustUploadFile([]byte("download me"), "download.txt", "text/plain", "")

	w := s.getFileDownloadLink(id)
	s.Equal(http.StatusOK, w.Code)
}

func (s *FileSuite) TestGetFileDownloadLink_ReturnsExpectedFields() {
	fileName := "fields_check.txt"
	id := s.mustUploadFile([]byte("content"), fileName, "text/plain", "")

	w := s.getFileDownloadLink(id)
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeFileDownloadLinkResponse(w)
	s.Equal(fileName, resp.FileName, "fileName should match the originally uploaded file name")
	s.NotEmpty(resp.DownloadURL, "downloadUrl should be non-empty")
	s.False(resp.ExpiresAt.IsZero(), "expiresAt should be a non-zero timestamp")
}

func (s *FileSuite) TestGetFileDownloadLink_ExpiresAtIsInFuture() {
	id := s.mustUploadFile([]byte("expiry check"), "expiry.txt", "text/plain", "")

	w := s.getFileDownloadLink(id)
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeFileDownloadLinkResponse(w)
	s.True(resp.ExpiresAt.After(time.Now()), "expiresAt must be a future timestamp")
}

// The usecase hardcodes a 15-minute presign duration (remove this later).
func (s *FileSuite) TestGetFileDownloadLink_ExpiresAtIsApproximately15MinutesFromNow() {
	id := s.mustUploadFile([]byte("timing check"), "timing.txt", "text/plain", "")

	before := time.Now()
	w := s.getFileDownloadLink(id)
	after := time.Now()
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeFileDownloadLinkResponse(w)

	lowerBound := before.Add(14 * time.Minute)
	upperBound := after.Add(16 * time.Minute)
	s.True(
		resp.ExpiresAt.After(lowerBound) && resp.ExpiresAt.Before(upperBound),
		"expiresAt %v should be roughly 15 minutes from now (window: %v – %v)",
		resp.ExpiresAt, lowerBound, upperBound,
	)
}

func (s *FileSuite) TestGetFileDownloadLink_NotFound_Returns404() {
	nonExistentID := "507f1f77bcf86cd799439011"

	w := s.getFileDownloadLink(nonExistentID)

	s.Equal(http.StatusNotFound, w.Code)
}

func (s *FileSuite) TestGetFileDownloadLink_InvalidIDFormat_Returns422() {
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
			w := s.getFileDownloadLink(tc.id)
			s.Equal(http.StatusUnprocessableEntity, w.Code, "case: %q", tc.name)
		})
	}
}

func (s *FileSuite) TestUploadThenGetDownloadLink_RoundTrip() {
	fileName := "roundtrip.txt"
	id := s.mustUploadFile([]byte("round trip content"), fileName, "text/plain", "")

	w := s.getFileDownloadLink(id)
	s.Require().Equal(http.StatusOK, w.Code)

	resp := s.decodeFileDownloadLinkResponse(w)
	s.Equal(fileName, resp.FileName)
	s.NotEmpty(resp.DownloadURL)
	s.True(resp.ExpiresAt.After(time.Now()))
}

func (s *FileSuite) TestUploadMultipleFiles_EachHasOwnDownloadLink() {
	type uploaded struct {
		id       string
		fileName string
	}

	files := []uploaded{
		{fileName: "alpha.txt"},
		{fileName: "beta.txt"},
		{fileName: "gamma.txt"},
	}

	for i := range files {
		buf := fmt.Appendf(nil, "content of %s", files[i].fileName)

		files[i].id = s.mustUploadFile(
			buf,
			files[i].fileName,
			"text/plain",
			"",
		)
	}

	downloadURLs := make(map[string]bool)
	for _, f := range files {
		w := s.getFileDownloadLink(f.id)
		s.Require().Equal(http.StatusOK, w.Code, "download link for %s should succeed", f.fileName)

		resp := s.decodeFileDownloadLinkResponse(w)
		s.Equal(f.fileName, resp.FileName, "fileName mismatch for %s", f.fileName)
		s.NotEmpty(resp.DownloadURL)
		downloadURLs[resp.DownloadURL] = true
	}

	s.Len(downloadURLs, len(files), "each file should have a distinct presigned download URL")
}
