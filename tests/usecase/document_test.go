package usecase_test

import (
	"bytes"
	"context"
	"docvault/entity"
	mock_test "docvault/tests/mock"
	"docvault/usecase"
	"errors"
	"io"
	"testing"
)

const TestPDF = "test.pdf"

func TestUploadSuccess(t *testing.T) {
	mockRepo := &mock_test.MockDocumentRepository{}
	mockStorage := &mock_test.MockServiceStorage{}
	mockQueue := &mock_test.MockServiceQueue{}

	mockRepo.SaveFunc = func(ctx context.Context, doc *entity.Document) error {
		return nil
	}

	mockStorage.UploadFunc = func(ctx context.Context, filename string, fileSize int64, contentType string, file io.Reader) error {
		return nil
	}

	mockQueue.PublishFunc = func(ctx context.Context, message string) error {
		return nil
	}

	uc := usecase.NewDocumentUsecase(mockRepo, mockStorage, mockQueue)

	doc, err := uc.Upload(context.Background(), TestPDF, 100, "application/pdf", bytes.NewReader([]byte("test content")), 60)
	if err != nil {
		t.Errorf("Upload() error = %v, want nil", err)
	}

	if doc == nil {
		t.Errorf("Upload() returned nil document, want non-nil")
	}

	if doc.FileName != TestPDF {
		t.Errorf("Upload() FileName = %s, want test.pdf", doc.FileName)
	}

	if doc.FileSize != 100 {
		t.Errorf("Upload() FileSize = %d, want 100", doc.FileSize)
	}
}

func TestUploadStorageFailure(t *testing.T) {
	mockRepo := &mock_test.MockDocumentRepository{}
	mockStorage := &mock_test.MockServiceStorage{}
	mockQueue := &mock_test.MockServiceQueue{}

	mockStorage.UploadFunc = func(ctx context.Context, filename string, fileSize int64, contentType string, file io.Reader) error {
		return errors.New("storage failed")
	}

	uc := usecase.NewDocumentUsecase(mockRepo, mockStorage, mockQueue)
	doc, err := uc.Upload(context.Background(), TestPDF, 100, "application/pdf", bytes.NewReader([]byte("test")), 60)

	if err == nil {
		t.Errorf("Upload() error = nil, want non-nil")
	}
	if doc != nil {
		t.Errorf("Upload() returned %v, want nil", doc)
	}
}

func TestDeleteSuccess(t *testing.T) {
	mockRepo := &mock_test.MockDocumentRepository{}
	mockStorage := &mock_test.MockServiceStorage{}
	mockQueue := &mock_test.MockServiceQueue{}

	mockRepo.FindByIdFunc = func(ctx context.Context, id string) (*entity.Document, error) {
		return &entity.Document{
			ID:       id,
			FileName: TestPDF,
		}, nil
	}

	mockQueue.PublishFunc = func(ctx context.Context, message string) error {
		return nil
	}

	mockRepo.DeleteFunc = func(ctx context.Context, id string) error {
		return nil
	}

	mockStorage.DeleteFunc = func(ctx context.Context, filename string) error {
		return nil
	}

	mockQueue.DeleteMessageFunc = func(ctx context.Context, receiptHandle string) error {
		return nil
	}

	uc := usecase.NewDocumentUsecase(mockRepo, mockStorage, mockQueue)

	if err := uc.Delete(context.Background(), "1"); err != nil {
		t.Errorf("Delete() error %v, want nil", err)
	}
}

func TestHealthAllServicesHealthy(t *testing.T) {
	mockRepo := &mock_test.MockDocumentRepository{}
	mockStorage := &mock_test.MockServiceStorage{}
	mockQueue := &mock_test.MockServiceQueue{}

	mockRepo.PingFunc = func(ctx context.Context) error {
		return nil
	}
	mockStorage.HealthFunc = func(ctx context.Context) error {
		return nil
	}
	mockQueue.HealthFunc = func(ctx context.Context) error {
		return nil
	}

	uc := usecase.NewDocumentUsecase(mockRepo, mockStorage, mockQueue)
	status := uc.Health(context.Background())

	if status["database"] != "ok" {
		t.Errorf("database status = %s, want ok", status["database"])
	}
	if status["storage"] != "ok" {
		t.Errorf("storage status = %s, want ok", status["storage"])
	}
	if status["queue"] != "ok" {
		t.Errorf("queue status = %s, want ok", status["queue"])
	}
}
