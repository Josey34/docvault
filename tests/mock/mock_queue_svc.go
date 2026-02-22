package mock_test

import "context"

type MockServiceQueue struct {
	PublishFunc       func(ctx context.Context, message string) error
	ConsumeFunc       func(ctx context.Context) (<-chan string, error)
	DeleteMessageFunc func(ctx context.Context, receiptHandle string) error
}

func (q *MockServiceQueue) Publish(ctx context.Context, message string) error {
	if q.PublishFunc != nil {
		return q.PublishFunc(ctx, message)
	}

	return nil
}

func (q *MockServiceQueue) Consume(ctx context.Context) (<-chan string, error) {
	if q.ConsumeFunc != nil {
		return q.ConsumeFunc(ctx)
	}

	return nil, nil
}

func (q *MockServiceQueue) DeleteMessage(ctx context.Context, receiptHandle string) error {
	if q.DeleteMessageFunc != nil {
		return q.DeleteMessageFunc(ctx, receiptHandle)
	}

	return nil
}
