package s3

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Downloader is the AWS S3 downloader wrapper
type Downloader struct {
	bucket     string
	downloader *s3manager.Downloader
}

// NewDownloader returns a new downloader
func NewDownloader(region, bucket string) *Downloader {
	downl := s3manager.NewDownloader(session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		SharedConfigState: session.SharedConfigEnable,
	})))

	return &Downloader{
		bucket:     bucket,
		downloader: downl,
	}
}

// NewDownloaderFromEnv returns a new downloader for provided env vars
func NewDownloaderFromEnv(regionEnv, bucketEnv string) *Downloader {
	region := os.Getenv(regionEnv)
	if region == "" {
		log.Fatalf("s3: environment variable '%s' is not set", regionEnv)
	}

	bucket := os.Getenv(bucketEnv)
	if bucket == "" {
		log.Fatalf("s3: environment variable '%s' is not set", bucketEnv)
	}

	return NewDownloader(region, bucket)
}

/*
DownloadFile downloads a file from an S3 bucket.
*/
func (d Downloader) DownloadFile(ctx context.Context, key string, body io.Writer) error {
	return d.DownloadFileFromBucket(ctx, key, d.bucket, body)
}

/*
DownloadFileFromBucket downloads a file from an S3 bucket.
*/
func (d Downloader) DownloadFileFromBucket(ctx context.Context, bucket, key string, body io.Writer) error {
	buff := new(aws.WriteAtBuffer)
	params := s3.GetObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
	}

	_, err := d.downloader.DownloadWithContext(ctx, buff, &params)
	if err != nil {
		return fmt.Errorf("failed to download file, bucket %s, key %s: %w", bucket, key, err)
	}

	_, err = body.Write(buff.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
