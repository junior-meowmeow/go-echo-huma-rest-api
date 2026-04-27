package s3api_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/storage/s3api"
)

func TestS3Storage(t *testing.T) {
	s3Client := setupS3Client(t)
	ctx := context.Background()
	bucketName := "test-bucket"

	_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	require.NoError(t, err)

	s3Storage := s3api.NewS3Storage(s3Client, bucketName)

	t.Run("UploadFile", func(t *testing.T) {
		cleanBucket(t, ctx, s3Client, bucketName)
		content := "Test content"
		reader := strings.NewReader(content)

		key := "test/data.txt"

		t.Run("Should upload file correctly", func(t *testing.T) {
			err := s3Storage.UploadFile(ctx, key, reader, int64(len(content)), "text/plain")

			require.NoError(t, err)

			// Verify content
			headObject, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
			})
			require.NoError(t, err)
			assert.Equal(t, int64(len(content)), *headObject.ContentLength)
		})
	})

	t.Run("CheckFileExists", func(t *testing.T) {
		cleanBucket(t, ctx, s3Client, bucketName)

		key := "exists.txt"
		_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   strings.NewReader("data"),
		})
		require.NoError(t, err)

		t.Run("Should return true when file exists", func(t *testing.T) {
			// Reuse uploaded file from previous step
			exists, err := s3Storage.CheckFileExists(ctx, key)

			require.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("Should return false when file not exists", func(t *testing.T) {
			exists, err := s3Storage.CheckFileExists(ctx, "non-existent-file.png")

			require.NoError(t, err)
			assert.False(t, exists)
		})
	})

	t.Run("ListFiles", func(t *testing.T) {
		cleanBucket(t, ctx, s3Client, bucketName)

		keysToUpload := []string{"a.txt", "b.txt"}
		for _, k := range keysToUpload {
			_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(k),
				Body:   strings.NewReader("data"),
			})
			require.NoError(t, err)
		}

		t.Run("Should success", func(t *testing.T) {
			keys, err := s3Storage.ListFiles(ctx, 10)

			require.NoError(t, err)
			assert.ElementsMatch(t, keys, keysToUpload)
		})

		t.Run("Should limit keys", func(t *testing.T) {
			keys, err := s3Storage.ListFiles(ctx, 1)
			require.NoError(t, err)
			assert.Len(t, keys, 1)
		})
	})

	t.Run("GetPresignedDownloadURL", func(t *testing.T) {
		cleanBucket(t, ctx, s3Client, bucketName)
		key := "privatefile.txt"
		filename := "download.txt"

		_, _ = s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   strings.NewReader("data"),
		})

		url, err := s3Storage.GetPresignedDownloadURL(ctx, key, filename, 10*time.Minute)

		t.Run("Should generate URL correctly", func(t *testing.T) {
			require.NoError(t, err)
			assert.NotEmpty(t, url)
			assert.Contains(t, url, bucketName)
			assert.Contains(t, url, filename)
			assert.Contains(t, url, "response-content-disposition=attachment")
			assert.Contains(t, url, "X-Amz-Expires=600")
			assert.Contains(t, url, "attachment%3B%20filename%3D%22download.txt%22")
		})
	})
}
