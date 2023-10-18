package s3

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Uploader is the AWS S3 uploader wrapper
type Uploader struct {
	bucket   string
	uploader *s3manager.Uploader
}

// NewUploader returns a new uploader
func NewUploader(region, bucket string) *Uploader {
	upl := s3manager.NewUploader(session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		SharedConfigState: session.SharedConfigEnable,
	})))

	return &Uploader{
		bucket:   bucket,
		uploader: upl,
	}
}

// NewUploaderFromEnv returns a new uploader for provided env vars
func NewUploaderFromEnv(regionEnv, bucketEnv string) *Uploader {
	region := os.Getenv(regionEnv)
	if region == "" {
		log.Fatalf("s3: environment variable '%s' is not set", regionEnv)
	}

	bucket := os.Getenv(bucketEnv)
	if bucket == "" {
		log.Fatalf("s3: environment variable '%s' is not set", bucketEnv)
	}

	return NewUploader(region, bucket)
}

/*
UploadFile uploads a file to an S3 bucket.
*/
func (u Uploader) UploadFile(ctx context.Context, key string, body io.Reader) error {
	params := s3manager.UploadInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
		Body:   body,
	}

	_, err := u.uploader.UploadWithContext(ctx, &params)
	if err != nil {
		return fmt.Errorf("upload file, bucket '%s', key '%s' fail: %w", u.bucket, key, err)
	}

	return nil
}
