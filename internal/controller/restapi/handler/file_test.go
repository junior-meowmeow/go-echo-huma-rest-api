package handler_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	usecasemocks "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase/mocks"
)

func TestFileHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("UploadFile", func(t *testing.T) {
		t.Run("Should upload file successfully", func(t *testing.T) {
			mockUC := usecasemocks.NewMockFileUseCase(t)

			filename := "test.txt"
			content := "hello world"
			contentType := "text/plain"
			baseKey := "my/folder"

			req := createMockUploadRequest(filename, content, contentType, baseKey)

			mockUC.EXPECT().
				UploadFile(
					mock.Anything,
					mock.Anything,
					filename,
					int64(len(content)),
					contentType,
					baseKey,
				).
				Return("file-id-123", nil)

			h := handler.NewFileHandler(mockUC)

			resp, err := h.UploadFile(context.Background(), req)

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, "file-id-123", resp.Body.ID)
		})

		t.Run("Should return error when upload fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockFileUseCase(t)

			req := createMockUploadRequest("test.txt", "content", "text/plain", "")

			mockUC.EXPECT().
				UploadFile(
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
				Return("", errors.New("upload error"))

			h := handler.NewFileHandler(mockUC)

			resp, err := h.UploadFile(context.Background(), req)

			require.Error(t, err)
			assert.Nil(t, resp)
			assert.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})

	t.Run("GetFileDownloadLink", func(t *testing.T) {
		t.Run("Should return download link", func(t *testing.T) {
			mockUC := usecasemocks.NewMockFileUseCase(t)

			now := time.Now()

			mockUC.EXPECT().
				GetFileDownloadLink(mock.Anything, "file-id").
				Return(entity.FileDownloadInfo{
					FileName:       "file.txt",
					DownloadURL:    "https://example.com",
					ExpirationTime: now,
				}, nil)

			h := handler.NewFileHandler(mockUC)

			req := &schema.GetFileDownloadLinkRequest{
				ID: "file-id",
			}

			resp, err := h.GetFileDownloadLink(ctx, req)

			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, "file.txt", resp.Body.FileName)
			assert.Equal(t, "https://example.com", resp.Body.DownloadURL)
			assert.Equal(t, now, resp.Body.ExpiresAt)
		})

		t.Run("Should return error when usecase fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockFileUseCase(t)

			mockUC.EXPECT().
				GetFileDownloadLink(mock.Anything, "file-id").
				Return(entity.FileDownloadInfo{}, errors.New("not found"))

			h := handler.NewFileHandler(mockUC)

			req := &schema.GetFileDownloadLinkRequest{
				ID: "file-id",
			}

			resp, err := h.GetFileDownloadLink(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)
			assert.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})

	t.Run("GetS3FileList", func(t *testing.T) {
		t.Run("Should return file list", func(t *testing.T) {
			mockUC := usecasemocks.NewMockFileUseCase(t)

			files := []string{"a.txt", "b.txt", "c.txt"}

			mockUC.EXPECT().
				GetS3FileList(mock.Anything).
				Return(files, nil)

			h := handler.NewFileHandler(mockUC)

			resp, err := h.GetS3FileList(ctx, &schema.GetS3FileListRequest{})

			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, files, resp.Body.Files)
			assert.Equal(t, 3, resp.Body.Count)
		})

		t.Run("Should return error when usecase fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockFileUseCase(t)

			mockUC.EXPECT().
				GetS3FileList(mock.Anything).
				Return(nil, errors.New("s3 error"))

			h := handler.NewFileHandler(mockUC)

			resp, err := h.GetS3FileList(ctx, &schema.GetS3FileListRequest{})

			require.Error(t, err)
			assert.Nil(t, resp)
			assert.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})
}

func createMockUploadRequest(filename, content, contentType, baseKey string) *schema.UploadFileRequest {
	form := &multipart.Form{
		Value: make(map[string][]string),
		File:  make(map[string][]*multipart.FileHeader),
	}

	form.Value["objectBaseKey"] = []string{baseKey}

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	h.Set("Content-Type", contentType)

	part, _ := mw.CreatePart(h)
	_, _ = part.Write([]byte(content))
	_ = mw.Close()

	reader := multipart.NewReader(&buf, mw.Boundary())
	parsedForm, _ := reader.ReadForm(1024)
	form.File["file"] = parsedForm.File["file"]

	fh := form.File["file"][0]
	file, _ := fh.Open()

	formFile := huma.FormFile{
		File:        file,
		ContentType: contentType,
		IsSet:       true,
		Size:        fh.Size,
		Filename:    fh.Filename,
	}

	mff := huma.MultipartFormFiles[schema.UploadFileData]{
		Form: form,
	}

	innerData := &schema.UploadFileData{
		File:          formFile,
		ObjectBaseKey: baseKey,
	}

	v := reflect.ValueOf(&mff).Elem()
	f := v.FieldByName("data")

	ptrToDataField := (**schema.UploadFileData)(unsafe.Pointer(f.UnsafeAddr()))
	*ptrToDataField = innerData

	return &schema.UploadFileRequest{
		RawBody: mff,
	}
}
