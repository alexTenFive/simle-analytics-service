package queue

import (
	"go.uber.org/zap"
	"testcase/server/internal/domain"
	"testcase/server/internal/entity/models"
	"testcase/server/internal/timestamp/usecase"
	"time"
)

const chTSBufferSize = 300

type ramQueue struct {
	repo    domain.Repository
	storeUC domain.InsertTimestampUseCase
	logger  *zap.SugaredLogger

	ticker     *time.Ticker
	bufferSize int
	cache      chan models.Timestamp
	tss        []models.Timestamp
	chExit     chan chan struct{}
}

func NewRamQueue(repo domain.Repository, exit chan chan struct{}, logger *zap.SugaredLogger) *ramQueue {
	rq := &ramQueue{
		storeUC:    usecase.NewInsertBatchUseCase(repo),
		cache:      make(chan models.Timestamp, chTSBufferSize),
		tss:        make([]models.Timestamp, 0, chTSBufferSize/2),
		bufferSize: chTSBufferSize,
		chExit:     exit,
		ticker:     time.NewTicker(time.Second),
		logger:     logger.With("service", "RamQueue"),
	}

	go rq.storeRoutine()

	return rq
}
func (x *ramQueue) Add(m models.Timestamp) {
	x.cache <- m
}

func (x *ramQueue) process() error {
	if err := x.storeUC.Do(x.tss); err != nil {
		return err
	}
	x.logger.Infof("saved: %d", len(x.tss))
	x.tss = x.tss[:0]

	return nil
}

func (x *ramQueue) storeRoutine() {
	for {
		select {
		case m, ok := <-x.cache:
			if !ok {
				continue
			}
			x.tss = append(x.tss, m)
			if len(x.tss) >= x.bufferSize/2 {
				if err := x.process(); err != nil {
					x.logger.Errorf("cannot store batch: %s", err)
					continue
				}
				x.ticker.Reset(time.Second)
			}
		case <-x.ticker.C:
			if len(x.tss) == 0 {
				continue
			}
			if err := x.process(); err != nil {
				x.logger.Errorf("cannot store batch: %s", err)
			}
		case exit := <-x.chExit:
			defer close(exit)
			x.logger.Warn("we are done here")
			if len(x.tss) == 0 {
				return
			}
			if err := x.process(); err != nil {
				x.logger.Errorf("cannot store batch: %s", err)
				return
			}
			return
		}
	}
}
