package usecase

import (
	"go.uber.org/zap"
	"testcase_v2/internal/entity/models"
	"time"
)

type timestampUseCase struct {
	repo   TimestampRepository
	queue  TimestampQueue
	logger *zap.SugaredLogger
}

func NewTimestampUseCase(repo TimestampRepository, queue TimestampQueue, logger *zap.SugaredLogger) *timestampUseCase {
	return &timestampUseCase{
		repo:   repo,
		queue:  queue,
		logger: logger.With("usecase", "Timestamp"),
	}
}

func (t *timestampUseCase) Insert(m models.Timestamp) error {
	t.queue.Add(m)
	return nil
}

func (t *timestampUseCase) GetAverageValueForLast(last time.Duration) (float32, error) {
	return t.repo.GetAverageValueForLast(last)
}
