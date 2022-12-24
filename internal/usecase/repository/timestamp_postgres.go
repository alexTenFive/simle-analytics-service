package repository

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"testcase_v2/internal/entity/models"
	"time"
)

type TimestampRepository struct {
	db     *sqlx.DB
	logger *zap.SugaredLogger
}

func NewDatabaseRepository(db *sqlx.DB, logger *zap.SugaredLogger) *TimestampRepository {
	return &TimestampRepository{
		db:     db,
		logger: logger.With("repository", "TimestampRepository"),
	}
}

func (x *TimestampRepository) InsertBatch(values []models.Timestamp) error {
	q := `INSERT INTO timestamps (ts, value) VALUES (:ts, :value)`
	if _, err := x.db.NamedExec(q, values); err != nil {
		return err
	}
	return nil
}

func (x *TimestampRepository) GetAverageValueForLast(last time.Duration) (float32, error) {
	var value float32
	if last == 0 {
		return 0, nil
	}

	q := `SELECT COALESCE(AVG("value"), 0) FROM timestamps WHERE "ts" BETWEEN $1 AND $2`
	if err := x.db.Get(&value, q, time.Now().Add(-last), time.Now()); err != nil {
		return -1, err
	}

	return value, nil
}
