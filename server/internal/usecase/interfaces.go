package usecase

import (
	"testcase_v2/server/internal/entity/models"
	"time"
)

type (
	TimestampUseCase interface {
		Insert(models.Timestamp) error
		GetAverageValueForLast(duration time.Duration) (float32, error)
	}
	TimestampRepository interface {
		InsertBatch([]models.Timestamp) error
		GetAverageValueForLast(time.Duration) (float32, error)
	}
	TimestampQueue interface {
		Add(v models.Timestamp)
	}
)
