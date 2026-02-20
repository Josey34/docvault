package usecase

import (
	"context"
	"docvault/entity"
	"docvault/repository"
	"docvault/service"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
)

type DocumentUsecase struct {
	repo    repository.DocumentRepository
	storage service.StorageService
}

func NewDocumentUsecase(repo repository.DocumentRepository, storage service.StorageService) *DocumentUsecase {
	return &DocumentUsecase{repo: repo, storage: storage}
}

func (u *DocumentUsecase) Upload(ctx context.Context, filename string, fileSize int64, contentType string, file io.Reader) (*entity.Document, error) {
	documentID := uuid.New().String()

	if err := u.storage.Upload(ctx, filename, fileSize, contentType, file); err != nil {
		return nil, fmt.Errorf("Failed to upload to storage %w", err)
	}

	document := &entity.Document{
		ID:          documentID,
		FileName:    filename,
		FileSize:    fileSize,
		ContentType: contentType,
		CreatedAt:   time.Now(),
		ExpiresAt:   nil,
	}
	if err := u.repo.Save(ctx, document); err != nil {
		return nil, fmt.Errorf("Failed to save to repository documents %w", err)
	}

	return document, nil
}

func (u *DocumentUsecase) List(ctx context.Context) ([]*entity.Document, error) {
	doc, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to list items from documents %w", err)
	}

	return doc, nil
}

func (u *DocumentUsecase) GetMetadata(ctx context.Context, id string) (*entity.Document, error) {
	doc, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to find item id %w", err)
	}

	return doc, nil
}
