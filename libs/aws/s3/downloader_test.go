package s3_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zale144/ube/libs/aws/s3"
)

func TestNewDownloader_Good(t *testing.T) {
	downloader := s3.NewDownloader("REGION", "BUCKET")
	require.NotNil(t, downloader)

	text := fmt.Sprintf("%#v", downloader)
	assert.Contains(t, text, `&s3.Downloader{bucket:"BUCKET"`)
}

func TestNewDownloaderFromEnv_Wrong_NoRegion(t *testing.T) {
	t.Setenv("REGION", "")
	t.Setenv("BUCKET", "myBucket")

	downloader, err := s3.NewDownloaderFromEnv("REGION", "BUCKET")
	require.Nil(t, downloader)
	require.Error(t, err)
}

func TestNewDownloaderFromEnv_Wrong_NoBucket(t *testing.T) {
	t.Setenv("REGION", "myRegion")
	t.Setenv("BUCKET", "")

	downloader, err := s3.NewDownloaderFromEnv("REGION", "BUCKET")
	require.Nil(t, downloader)
	require.Error(t, err)
}

func TestNewDownloaderFromEnv_Good(t *testing.T) {
	t.Setenv("REGION", "myRegion")
	t.Setenv("BUCKET", "myBucket")

	downloader, err := s3.NewDownloaderFromEnv("REGION", "BUCKET")
	require.NotNil(t, downloader)
	require.NoError(t, err)

	text := fmt.Sprintf("%#v", downloader)
	assert.Contains(t, text, `&s3.Downloader{bucket:"myBucket"`)
}

func TestDownloadFile_Wrong(t *testing.T) {
	downloader := s3.NewDownloader("REGION", "BUCKET")
	require.NotNil(t, downloader)

	ctx := context.Background()
	var body bytes.Buffer
	err := downloader.DownloadFile(ctx, "testkey", &body)
	require.Error(t, err)
}
