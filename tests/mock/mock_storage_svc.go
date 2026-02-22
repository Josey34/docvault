package mock_test

import (
	"context"
	"io"
)

type MockServiceStorage struct {
	UploadFunc   func(ctx context.Context, filename string, fileSize int64, contentType string, file io.Reader) error
	DownloadFunc func(ctx context.Context, filename string) (io.ReadCloser, error)
	DeleteFunc   func(ctx context.Context, filename string) error
}

func (s *MockServiceStorage) Upload(ctx context.Context, filename string, fileSize int64, contentType string, file io.Reader) error {
	if s.UploadFunc != nil {
		return s.UploadFunc(ctx, filename, fileSize, contentType, file)
	}

	return nil
}

func (s *MockServiceStorage) Download(ctx context.Context, filename string) (io.ReadCloser, error) {
	if s.DownloadFunc != nil {
		return s.DownloadFunc(ctx, filename)
	}

	return nil, nil
}

func (s *MockServiceStorage) Delete(ctx context.Context, filename string) error {
	if s.DeleteFunc != nil {
		return s.DeleteFunc(ctx, filename)
	}

	return nil
}
