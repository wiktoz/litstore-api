package upload

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type R2Uploader struct {
	client   *s3.Client
	bucket   string
	endpoint string
}

func NewR2Uploader(client *s3.Client, bucket string, endpoint string) *R2Uploader {
	return &R2Uploader{client: client, bucket: bucket, endpoint: endpoint}
}

func (r *R2Uploader) Upload(ctx context.Context, filename string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	_, err = r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:             aws.String(r.bucket),
		Key:                aws.String(filename),
		Body:               file,
		ContentType:        aws.String(fileHeader.Header.Get("Content-Type")),
		ContentDisposition: aws.String("inline"),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
			return "", fmt.Errorf("object is too large. Use multipart upload or console")
		}
		return "", fmt.Errorf("upload failed: %w", err)
	}

	err = s3.NewObjectExistsWaiter(r.client).Wait(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(filename),
	}, time.Minute)

	if err != nil {
		return "", fmt.Errorf("failed to confirm object exists after upload: %w", err)
	}

	publicURL := fmt.Sprintf("%s/%s", "https://assets.litstore.pl", filename)
	return publicURL, nil
}
