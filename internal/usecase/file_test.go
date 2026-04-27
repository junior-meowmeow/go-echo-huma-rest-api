package usecase_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/port/mocks"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

func TestFileUseCase(t *testing.T) {
	ctx := context.Background()

	t.Run("UploadFile", func(t *testing.T) {
		t.Run("Should upload and save file record successfully", func(t *testing.T) {
			mockStorage := mocks.NewMockFileStorage(t)
			mockRepo := mocks.NewMockFileRecordRepository(t)

			fileContent := "file-content"
			fileSize := int64(len(fileContent))
			stream := bytes.NewBufferString(fileContent)

			mockStorage.EXPECT().
				CheckFileExists(ctx, mock.Anything).
				Return(false, nil)

			mockStorage.EXPECT().
				UploadFile(ctx, mock.Anything, stream, fileSize, "text/plain").
				Return(nil)

			mockRepo.EXPECT().
				CreateFileRecord(ctx, mock.Anything).
				Run(func(_ context.Context, rec *entity.FileRecord) {
					assert.Equal(t, "test.txt", rec.FileName)
					assert.Equal(t, fileSize, rec.Size)
					assert.Equal(t, "text/plain", rec.ContentType)
					assert.NotEmpty(t, rec.S3Key)

					assert.False(t, rec.CreatedAt.IsZero())
					assert.False(t, rec.UpdatedAt.IsZero())
					assert.WithinDuration(t, rec.CreatedAt, rec.UpdatedAt, time.Millisecond)
				}).
				Return("file-id", nil)

			uc := usecase.NewFileUseCase(mockRepo, mockStorage)

			id, err := uc.UploadFile(ctx, stream, "test.txt", fileSize, "text/plain", "base/")

			require.NoError(t, err)
			assert.Equal(t, "file-id", id)
		})

		t.Run("Should return error when CheckFileExists fails", func(t *testing.T) {
			mockStorage := mocks.NewMockFileStorage(t)
			mockRepo := mocks.NewMockFileRecordRepository(t)

			mockErr := errors.New("s3 error")

			mockStorage.EXPECT().
				CheckFileExists(ctx, mock.Anything).
				Return(false, mockErr)

			uc := usecase.NewFileUseCase(mockRepo, mockStorage)

			_, err := uc.UploadFile(ctx, bytes.NewBuffer(nil), "a.txt", 1, "text/plain", "")

			require.Error(t, err)
			require.ErrorContains(t, err, "failed to check file existence")
			assert.ErrorIs(t, err, mockErr)
		})

		t.Run("Should return error when upload fails", func(t *testing.T) {
			mockStorage := mocks.NewMockFileStorage(t)
			mockRepo := mocks.NewMockFileRecordRepository(t)

			mockErr := errors.New("upload error")

			mockStorage.EXPECT().
				CheckFileExists(ctx, mock.Anything).
				Return(false, nil)

			mockStorage.EXPECT().
				UploadFile(ctx, mock.Anything, mock.Anything, int64(1), "text/plain").
				Return(mockErr)

			uc := usecase.NewFileUseCase(mockRepo, mockStorage)

			_, err := uc.UploadFile(ctx, bytes.NewBuffer(nil), "a.txt", 1, "text/plain", "")

			require.Error(t, err)
			require.ErrorContains(t, err, "failed to upload to S3")
			assert.ErrorIs(t, err, mockErr)
		})

		t.Run("Should return error when saving file record fails", func(t *testing.T) {
			mockStorage := mocks.NewMockFileStorage(t)
			mockRepo := mocks.NewMockFileRecordRepository(t)

			mockErr := errors.New("db error")

			mockStorage.EXPECT().
				CheckFileExists(ctx, mock.Anything).
				Return(false, nil)

			mockStorage.EXPECT().
				UploadFile(ctx, mock.Anything, mock.Anything, int64(1), "text/plain").
				Return(nil)

			mockRepo.EXPECT().
				CreateFileRecord(ctx, mock.Anything).
				Return("", mockErr)

			uc := usecase.NewFileUseCase(mockRepo, mockStorage)

			_, err := uc.UploadFile(ctx, bytes.NewBuffer(nil), "a.txt", 1, "text/plain", "")

			require.Error(t, err)
			require.ErrorContains(t, err, "failed to save file record")
			assert.ErrorIs(t, err, mockErr)
		})
	})

	t.Run("GetFileDownloadLink", func(t *testing.T) {
		t.Run("Should return download link", func(t *testing.T) {
			mockStorage := mocks.NewMockFileStorage(t)
			mockRepo := mocks.NewMockFileRecordRepository(t)

			record := entity.FileRecord{
				ID:       "1",
				FileName: "file.txt",
				S3Key:    "key-123",
			}

			mockRepo.EXPECT().
				GetFileRecordByID(ctx, "1").
				Return(record, nil)

			mockStorage.EXPECT().
				GetPresignedDownloadURL(ctx, "key-123", "file.txt", 15*time.Minute).
				Return("http://signed-url", nil)

			uc := usecase.NewFileUseCase(mockRepo, mockStorage)

			result, err := uc.GetFileDownloadLink(ctx, "1")

			require.NoError(t, err)
			assert.Equal(t, "http://signed-url", result.DownloadURL)
			assert.Equal(t, "file.txt", result.FileName)
			assert.False(t, result.ExpirationTime.IsZero())
		})

		t.Run("Should return error when record not found", func(t *testing.T) {
			mockStorage := mocks.NewMockFileStorage(t)
			mockRepo := mocks.NewMockFileRecordRepository(t)

			mockErr := errors.New("not found")

			mockRepo.EXPECT().
				GetFileRecordByID(ctx, "1").
				Return(entity.FileRecord{}, mockErr)

			uc := usecase.NewFileUseCase(mockRepo, mockStorage)

			_, err := uc.GetFileDownloadLink(ctx, "1")

			require.Error(t, err)
			require.ErrorContains(t, err, "file not found")
			assert.ErrorIs(t, err, mockErr)
		})

		t.Run("Should return error when presign fails", func(t *testing.T) {
			mockStorage := mocks.NewMockFileStorage(t)
			mockRepo := mocks.NewMockFileRecordRepository(t)

			record := entity.FileRecord{
				ID:       "1",
				FileName: "file.txt",
				S3Key:    "key-123",
			}

			mockErr := errors.New("presign error")

			mockRepo.EXPECT().
				GetFileRecordByID(ctx, "1").
				Return(record, nil)

			mockStorage.EXPECT().
				GetPresignedDownloadURL(ctx, "key-123", "file.txt", 15*time.Minute).
				Return("", mockErr)

			uc := usecase.NewFileUseCase(mockRepo, mockStorage)

			_, err := uc.GetFileDownloadLink(ctx, "1")

			require.Error(t, err)
			require.ErrorContains(t, err, "failed to presign url")
			assert.ErrorIs(t, err, mockErr)
		})
	})

	t.Run("GetS3FileList", func(t *testing.T) {
		t.Run("Should return file list", func(t *testing.T) {
			mockStorage := mocks.NewMockFileStorage(t)
			mockRepo := mocks.NewMockFileRecordRepository(t)

			expected := []string{"a.txt", "b.txt"}

			mockStorage.EXPECT().
				ListFiles(ctx, 20).
				Return(expected, nil)

			uc := usecase.NewFileUseCase(mockRepo, mockStorage)

			result, err := uc.GetS3FileList(ctx)

			require.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("Should return wrapped error", func(t *testing.T) {
			mockStorage := mocks.NewMockFileStorage(t)
			mockRepo := mocks.NewMockFileRecordRepository(t)

			mockErr := errors.New("s3 error")

			mockStorage.EXPECT().
				ListFiles(ctx, 20).
				Return(nil, mockErr)

			uc := usecase.NewFileUseCase(mockRepo, mockStorage)

			result, err := uc.GetS3FileList(ctx)

			require.Error(t, err)
			assert.Nil(t, result)
			require.ErrorContains(t, err, "failed to list S3 files")
			assert.ErrorIs(t, err, mockErr)
		})
	})
}
