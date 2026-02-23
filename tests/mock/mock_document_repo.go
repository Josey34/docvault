package mock_test

import (
	"context"
	"docvault/entity"
	"time"
)

type MockDocumentRepository struct {
	SaveFunc        func(ctx context.Context, doc *entity.Document) error
	FindByIdFunc    func(ctx context.Context, id string) (*entity.Document, error)
	FindAllFunc     func(ctx context.Context) ([]*entity.Document, error)
	DeleteFunc      func(ctx context.Context, id string) error
	FindExpiredFunc func(ctx context.Context, now time.Time) ([]*entity.Document, error)
	PingFunc        func(ctx context.Context) error
}

func (m *MockDocumentRepository) Save(ctx context.Context, doc *entity.Document) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, doc)
	}

	return nil
}

func (m *MockDocumentRepository) FindById(ctx context.Context, id string) (*entity.Document, error) {
	if m.FindByIdFunc != nil {
		return m.FindByIdFunc(ctx, id)
	}

	return nil, nil
}

func (m *MockDocumentRepository) FindAll(ctx context.Context) ([]*entity.Document, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}

	return nil, nil
}

func (m *MockDocumentRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}

	return nil
}

func (m *MockDocumentRepository) FindExpired(ctx context.Context, now time.Time) ([]*entity.Document, error) {
	if m.FindExpiredFunc != nil {
		return m.FindExpiredFunc(ctx, now)
	}

	return nil, nil
}

func (m *MockDocumentRepository) Ping(ctx context.Context) error {
	if m.PingFunc != nil {
		return m.PingFunc(ctx)
	}

	return nil
}
