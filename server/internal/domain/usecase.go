package domain

import (
	"testcase/server/internal/entity/models"
	"time"
)

type InsertTimestampUseCase interface {
	Do([]models.Timestamp) error
}

type GetAverageValueUseCase interface {
	Do(time.Duration) (float32, error)
}
