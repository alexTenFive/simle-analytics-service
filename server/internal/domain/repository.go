package domain

import (
	"testcase/server/internal/entity/models"
	"time"
)

type Repository interface {
	InsertBatchTS([]models.Timestamp) error
	GetAverageValueTSForLast(time.Duration) (float32, error)
}
