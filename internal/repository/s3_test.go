package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS3Repository(t *testing.T) {
	s3Client := setupS3Client(t)
	ctx := context.Background()

	bucketName := "test-bucket"
	_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	require.NoError(t, err)

	repo := repository.NewS3Repository(s3Client, bucketName)

	t.Run("UploadFile", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-file-*.txt")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		content := "Repository Integration Test Content"
		_, err = tmpFile.WriteString(content)
		require.NoError(t, err)

		// Reset cursor to start of file
		_, err = tmpFile.Seek(0, 0)
		require.NoError(t, err)

		key := "folder/test.txt"

		err = repo.UploadFile(ctx, key, tmpFile, int64(len(content)), "text/plain")
		assert.NoError(t, err)

		// Verification: Check if it really exists in S3
		_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		assert.NoError(t, err, "File should exist in S3")
	})

	t.Run("ListFiles", func(t *testing.T) {
		// We expect the file from the previous test ("folder/test.txt")
		keys, err := repo.ListFiles(ctx, 10)
		assert.NoError(t, err)

		require.Len(t, keys, 1)
		assert.Equal(t, "folder/test.txt", keys[0])
	})

	t.Run("GetPresignedDownloadURL", func(t *testing.T) {
		key := "folder/test.txt"
		filename := "download-me.txt"

		url, err := repo.GetPresignedDownloadURL(ctx, key, filename, 15*time.Minute)
		assert.NoError(t, err)
		assert.NotEmpty(t, url)

		t.Logf("Generated Presigned URL: %s", url)
		assert.Contains(t, url, bucketName)
		assert.Contains(t, url, "folder/test.txt")
		assert.Contains(t, url, "response-content-disposition=attachment")
	})
}
