package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Client struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

func NewR2Client(accountID, accessKeyID, secretKey, bucket, publicURL string) *R2Client {
	r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	client := s3.New(s3.Options{
		BaseEndpoint: aws.String(r2Endpoint),
		Region:       "auto",
		Credentials:  credentials.NewStaticCredentialsProvider(accessKeyID, secretKey, ""),
	})

	return &R2Client{
		client:    client,
		bucket:    bucket,
		publicURL: publicURL,
	}
}

// UploadFile takes an open multipart file, uploads it to R2,
// and returns the public URL you can store in the database.
func (r *R2Client) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate file type — only images allowed
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
	}
	contentType, ok := allowed[ext]
	if !ok {
		return "", fmt.Errorf("file type %s not allowed", ext)
	}

	// Unique filename: timestamp + original name, no spaces
	key := fmt.Sprintf("%d-%s", time.Now().UnixNano(), sanitiseFilename(header.Filename))

	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}

	return fmt.Sprintf("%s/%s", r.publicURL, key), nil
}

func sanitiseFilename(name string) string {
	base := filepath.Base(name)
	// Replace spaces with hyphens, keep it URL-safe
	return strings.ReplaceAll(base, " ", "-")
}
