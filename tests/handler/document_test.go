package handler_test

import (
	"bytes"
	"context"
	"docvault/dto"
	"docvault/entity"
	"docvault/handler"
	mock_test "docvault/tests/mock"
	"docvault/usecase"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

const (
	TestPDFFileName = "test.pdf"
	TestPDFContent  = "pdf content"
	TestExpiresIn   = "60"
)

func createDefaultMocks(errSaveFunc error, errUploadFunc error, errPublishFunc error) (*mock_test.MockDocumentRepository, *mock_test.MockServiceStorage, *mock_test.MockServiceQueue) {
	mockRepo := &mock_test.MockDocumentRepository{}
	mockStorage := &mock_test.MockServiceStorage{}
	mockQueue := &mock_test.MockServiceQueue{}

	mockRepo.SaveFunc = func(ctx context.Context, doc *entity.Document) error {
		return errSaveFunc
	}

	mockStorage.UploadFunc = func(ctx context.Context, filename string, fileSize int64, contentType string, file io.Reader) error {
		return errUploadFunc
	}

	mockQueue.PublishFunc = func(ctx context.Context, message string) error {
		return errPublishFunc
	}

	return mockRepo, mockStorage, mockQueue
}

func TestUploadHandlerSuccess(t *testing.T) {
	mockRepo, mockStorage, mockQueue := createDefaultMocks(nil, nil, nil)

	uc := usecase.NewDocumentUsecase(mockRepo, mockStorage, mockQueue)
	h := handler.NewDocumentHandler(uc)

	router := gin.New()
	router.POST("/upload", h.Upload)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, _ := writer.CreateFormFile("file", TestPDFFileName)
	file.Write([]byte(TestPDFContent))
	writer.WriteField("expires_in", TestExpiresIn)
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var response dto.DocumentResponse

	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("Upload() status = %d, want %d", rec.Code, http.StatusCreated)
	}

	if response.FileName == "" {
		t.Errorf("Upload() FileName = %s, want non-empty", response.FileName)
	}

	if response.FileSize == 0 {
		t.Errorf("Upload() FileSize = %d, want non-empty", response.FileSize)
	}

	if response.ID == "" {
		t.Errorf("Upload() ID = %s, want non-empty", response.ID)
	}
}

func TestUploadHandlerBadRequest(t *testing.T) {
	mockRepo, mockStorage, mockQueue := createDefaultMocks(nil, nil, nil)

	uc := usecase.NewDocumentUsecase(mockRepo, mockStorage, mockQueue)
	h := handler.NewDocumentHandler(uc)

	router := gin.New()
	router.POST("/upload", h.Upload)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("expires_in", TestExpiresIn)
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Upload() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] == "" {
		t.Errorf("Upload() error = %s, want non-empty", response["error"])
	}
}

func TestListHandlerSuccess(t *testing.T) {
	mockRepo, mockStorage, mockQueue := createDefaultMocks(nil, nil, nil)

	mockRepo.FindAllFunc = func(ctx context.Context) ([]*entity.Document, error) {
		return []*entity.Document{
			{
				ID:       "test-id-1",
				FileName: "doc1.pdf",
				FileSize: 1024,
			},
			{
				ID:       "test-id-2",
				FileName: "doc2.pdf",
				FileSize: 2048,
			},
		}, nil
	}

	uc := usecase.NewDocumentUsecase(mockRepo, mockStorage, mockQueue)
	h := handler.NewDocumentHandler(uc)

	router := gin.New()
	router.GET("/api/documents", h.List)

	req := httptest.NewRequest("GET", "/api/documents", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("List() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var response []*dto.DocumentResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) == 0 {
		t.Errorf("List() returned empty slice, want non-empty")
	}

	if response[0].FileName != "doc1.pdf" {
		t.Errorf("List() first doc FileName = %s, want doc1.pdf", response[0].FileName)
	}
}

func TestDeleteHandlerSuccess(t *testing.T) {
	mockRepo, mockStorage, mockQueue := createDefaultMocks(nil, nil, nil)

	mockRepo.FindByIdFunc = func(ctx context.Context, id string) (*entity.Document, error) {
		return &entity.Document{
			ID:       id,
			FileName: "doc1.pdf",
			FileSize: 100,
		}, nil
	}

	mockRepo.DeleteFunc = func(ctx context.Context, id string) error {
		return nil
	}

	uc := usecase.NewDocumentUsecase(mockRepo, mockStorage, mockQueue)
	h := handler.NewDocumentHandler(uc)

	router := gin.New()
	router.DELETE("/api/documents/:id", h.Delete)

	req := httptest.NewRequest("DELETE", "/api/documents/test-id-123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Delete() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "document deleted" {
		t.Errorf("Delete() message = %s, want 'document deleted'", response["message"])
	}
}

func TestDownloadHandlerSuccess(t *testing.T) {
	mockRepo, mockStorage, mockQueue := createDefaultMocks(nil, nil, nil)

	mockRepo.FindByIdFunc = func(ctx context.Context, id string) (*entity.Document, error) {
		return &entity.Document{
			ID:          id,
			FileName:    "doc1.pdf",
			FileSize:    100,
			ContentType: "application/pdf",
		}, nil
	}

	mockStorage.DownloadFunc = func(ctx context.Context, filename string) (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader([]byte(TestPDFContent))), nil
	}

	uc := usecase.NewDocumentUsecase(mockRepo, mockStorage, mockQueue)
	h := handler.NewDocumentHandler(uc)

	router := gin.New()
	router.GET("/api/documents/:id/download", h.Download)

	req := httptest.NewRequest("GET", "/api/documents/test-id-123/download", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GET() status = %d, want %d", rec.Code, http.StatusOK)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/pdf" {
		t.Errorf("GET() Content-Type = %s, want application/pdf", contentType)
	}

	contentDisposition := rec.Header().Get("Content-Disposition")
	expectedDisposition := "attachment; filename=\"doc1.pdf\""
	if contentDisposition != expectedDisposition {
		t.Errorf("GET() Content-Disposition = %s, want %s", contentDisposition, expectedDisposition)
	}

	bodyBytes := rec.Body.Bytes()
	if !bytes.Equal(bodyBytes, []byte(TestPDFContent)) {
		t.Errorf("GET() body = %s, want %s", string(bodyBytes), TestPDFContent)
	}
}

func TestHealthHandlerSuccess(t *testing.T) {
	mockRepo, mockStorage, mockQueue := createDefaultMocks(nil, nil, nil)

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
	h := handler.NewDocumentHandler(uc)

	router := gin.New()
	router.GET("/health", h.Health)

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Health() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Health() status = %v, want 'healthy'", response["status"])
	}

	services := response["services"].(map[string]interface{})
	if services["database"] != "ok" {
		t.Errorf("Health() database = %v, want 'ok'", services["database"])
	}
	if services["storage"] != "ok" {
		t.Errorf("Health() storage = %v, want 'ok'", services["storage"])
	}
	if services["queue"] != "ok" {
		t.Errorf("Health() queue = %v, want 'ok'", services["queue"])
	}
}
