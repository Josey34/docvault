package service

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinIOStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinIOStorage(client *minio.Client, bucketName string) StorageService {
	return &MinIOStorage{
		client:     client,
		bucketName: bucketName,
	}
}

func (m *MinIOStorage) Upload(ctx context.Context, filename string, fileSize int64, contentType string, file io.Reader) error {
	info, err := m.client.PutObject(ctx, m.bucketName, filename, file, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("Error initializing minio upload %w", err)
	}

	fmt.Printf("Successfully uploaded %s of size %d\n", info.Key, info.Size)

	return nil
}

func (m *MinIOStorage) Download(ctx context.Context, filename string) (io.ReadCloser, error) {
	object, err := m.client.GetObject(ctx, m.bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("Error initialize minio download %w", err)
	}

	return object, nil
}

func (m *MinIOStorage) Delete(ctx context.Context, filename string) error {
	if err := m.client.RemoveObject(ctx, m.bucketName, filename, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("Error deleting object from minio %w", err)
	}

	return nil
}
