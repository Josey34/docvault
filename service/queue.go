package service

import "context"

type QueueService interface {
	Publish(ctx context.Context, message string) error
	Consume(ctx context.Context) (<-chan string, error)
	DeleteMessage(ctx context.Context, receiptHandle string) error
}
