package queue

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"testcase_v2/internal/entity/models"
	"testcase_v2/internal/usecase"
	"time"
)

const chTSBufferSize = 60

type ramQueue struct {
	repo   usecase.TimestampRepository
	logger *zap.SugaredLogger

	ticker     *time.Ticker
	bufferSize int
	cache      chan models.Timestamp
	tss        []models.Timestamp
	chExit     chan chan struct{}
}

func NewRamQueue(repo usecase.TimestampRepository, exit chan chan struct{}, logger *zap.SugaredLogger) *ramQueue {
	bufferSize := chTSBufferSize
	cBufferSize := viper.GetInt("server.queue_buffer_size")
	if cBufferSize > 0 && cBufferSize < 1e5 {
		logger.Debugf("get buffer size from config: %d", cBufferSize)
		bufferSize = cBufferSize
	}
	rq := &ramQueue{
		repo:       repo,
		cache:      make(chan models.Timestamp, bufferSize),
		tss:        make([]models.Timestamp, 0, bufferSize),
		bufferSize: bufferSize,
		chExit:     exit,
		ticker:     time.NewTicker(time.Second),
		logger:     logger.With("service", "timestamp_ram_queue"),
	}

	go rq.storeRoutine()

	return rq
}
func (x *ramQueue) Add(m models.Timestamp) {
	x.cache <- m
}

func (x *ramQueue) process() error {
	if err := x.repo.InsertBatch(x.tss); err != nil {
		return err
	}
	x.logger.Debugf("saved: %d", len(x.tss))
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
			if len(x.tss) >= x.bufferSize {
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
			x.logger.Debug("we are done here")
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
