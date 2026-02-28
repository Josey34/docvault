package worker

import (
	"context"
	"docvault/service"
	"fmt"
	"log"
)

type NotificationWorker struct {
	queue service.QueueService
}

func NewNotificationWorker(queue service.QueueService) *NotificationWorker {
	return &NotificationWorker{
		queue: queue,
	}
}

func (w *NotificationWorker) Start(ctx context.Context) {
	msgChan, err := w.queue.Consume(ctx)
	if err != nil {
		log.Println("Error consuming queue:", err)
		return
	}

	for message := range msgChan {
		fmt.Printf("Received message: %s\n", message)
	}
}
