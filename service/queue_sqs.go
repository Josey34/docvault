package service

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSQueue struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSQueue(client *sqs.Client, queueURL string) QueueService {
	return &SQSQueue{
		client:   client,
		queueURL: queueURL,
	}
}

func (s *SQSQueue) Publish(ctx context.Context, message string) error {
	_, err := s.client.SendMessage(ctx, &sqs.SendMessageInput{MessageBody: &message, QueueUrl: &s.queueURL})
	if err != nil {
		return fmt.Errorf("Failed to publish sqs queue %w", err)
	}

	return nil
}

func (s *SQSQueue) Consume(ctx context.Context) (<-chan string, error) {
	messageChan := make(chan string, 10)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(messageChan)
				return
			default:
				queueUrl := aws.String(s.queueURL)
				output, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
					QueueUrl:            queueUrl,
					MaxNumberOfMessages: *aws.Int32(10),
					WaitTimeSeconds:     20,
				})
				if err != nil {
					fmt.Printf("Error receiving message from SQS: %v\n", err)
					continue
				}
				for _, msg := range output.Messages {
					messageChan <- *msg.Body
				}
			}
		}
	}()
	return messageChan, nil
}

func (s *SQSQueue) DeleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{QueueUrl: &s.queueURL, ReceiptHandle: &receiptHandle})
	if err != nil {
		return fmt.Errorf("Error delete message %w", err)
	}

	return nil
}
