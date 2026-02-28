package worker

import (
	"context"
	"docvault/usecase"
	"time"
)

type SchedulerWorker struct {
	usecase *usecase.DocumentUsecase
}

func NewSchedulerWorker(usecase *usecase.DocumentUsecase) *SchedulerWorker {
	return &SchedulerWorker{usecase: usecase}
}

func (s *SchedulerWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-ticker.C:
			s.usecase.DeleteExpiredDocuments(ctx)
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
