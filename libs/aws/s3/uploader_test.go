package s3_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zale144/ube/libs/aws/s3"
)

func TestNewUploader_Good(t *testing.T) {
	uploader := s3.NewUploader("REGION", "BUCKET")
	require.NotNil(t, uploader)

	text := fmt.Sprintf("%#v", uploader)
	assert.Contains(t, text, `&s3.Uploader{bucket:"BUCKET"`)
}

func TestNewUploaderFromEnv_Wrong_NoRegion(t *testing.T) {
	t.Setenv("REGION", "")
	t.Setenv("BUCKET", "myBucket")

	uploader, err := s3.NewUploaderFromEnv("REGION", "BUCKET")
	require.Nil(t, uploader)
	require.Error(t, err)
}

func TestNewUploaderFromEnv_Wrong_NoBucket(t *testing.T) {
	t.Setenv("REGION", "myRegion")
	t.Setenv("BUCKET", "")

	uploader, err := s3.NewUploaderFromEnv("REGION", "BUCKET")
	require.Nil(t, uploader)
	require.Error(t, err)
}

func TestNewUploaderFromEnv_Good(t *testing.T) {
	t.Setenv("REGION", "myRegion")
	t.Setenv("BUCKET", "myBucket")

	uploader, err := s3.NewUploaderFromEnv("REGION", "BUCKET")
	require.NotNil(t, uploader)
	require.NoError(t, err)

	text := fmt.Sprintf("%#v", uploader)
	assert.Contains(t, text, `&s3.Uploader{bucket:"myBucket"`)
}

func TestUploadFile_Wrong(t *testing.T) {
	uploader := s3.NewUploader("REGION", "BUCKET")
	require.NotNil(t, uploader)

	ctx := context.Background()
	body := strings.NewReader("file content")

	err := uploader.UploadFile(ctx, "testkey", body)
	require.Error(t, err)
}
