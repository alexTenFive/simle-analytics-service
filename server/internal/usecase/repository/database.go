package repository

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"testcase/server/internal/entity/models"
	"time"
)

type DatabaseRepository struct {
	db     *sqlx.DB
	logger *zap.SugaredLogger
}

func NewDatabaseRepository(db *sqlx.DB, logger *zap.SugaredLogger) *DatabaseRepository {
	return &DatabaseRepository{
		db:     db,
		logger: logger.With("repository", "DatabaseRepository"),
	}
}

func (dr *DatabaseRepository) InsertBatchTS(values []models.Timestamp) error {
	q := `INSERT INTO timestamps (ts, value) VALUES (:ts, :value)`
	if _, err := dr.db.NamedExec(q, values); err != nil {
		return err
	}
	return nil
}

func (dr *DatabaseRepository) GetAverageValueTSForLast(last time.Duration) (float32, error) {
	var value float32
	if last == 0 {
		return 0, nil
	}

	q := `SELECT AVG("value") FROM timestamps WHERE "ts" BETWEEN $1 AND $2`
	if err := dr.db.Get(&value, q, time.Now().Add(-last), time.Now()); err != nil {
		return -1, err
	}

	return value, nil
}
