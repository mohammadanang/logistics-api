package gcp

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"cloud.google.com/go/storage"
)

type StorageClient struct {
	client *storage.Client
	bucket string
}

func NewStorageClient(bucketName string) (*StorageClient, error) {
	ctx := context.Background()
	// Secara otomatis akan membaca file JSON dari environment variable GOOGLE_APPLICATION_CREDENTIALS
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCP storage client: %v", err)
	}
	return &StorageClient{client: client, bucket: bucketName}, nil
}

// UploadFile mengunggah file ke bucket dan mengembalikan URL publiknya
func (s *StorageClient) UploadFile(ctx context.Context, file multipart.File, objectName string) (string, error) {
	wc := s.client.Bucket(s.bucket).Object(objectName).NewWriter(ctx)

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to copy file to bucket: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Format URL Publik GCS
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucket, objectName)
	return publicURL, nil
}
