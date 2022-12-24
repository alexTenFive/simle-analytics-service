package http

import (
	"go.uber.org/zap"
	"net/http"
	"testcase_v2/internal/entity/models"
	"testcase_v2/internal/entity/transport"
	"testcase_v2/internal/usecase"
	"testcase_v2/pkg/httpErrors"
	"testcase_v2/pkg/response"
	"time"
)

type timestampHandlers struct {
	ts     usecase.TimestampUseCase
	logger *zap.SugaredLogger
}

func newTimestampHandlers(ts usecase.TimestampUseCase, logger *zap.SugaredLogger) *timestampHandlers {
	return &timestampHandlers{
		ts:     ts,
		logger: logger.With("handlers", "timestamps"),
	}
}

func (t *timestampHandlers) Receive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		item := transport.SendRequest{}
		if err := item.Decode(r.Body); err != nil {
			t.logger.Errorf("decode request:%s", err)
			http.Error(w, httpErrors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}
		if item.Timestamp == 0 {
			http.Error(w, httpErrors.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		if err := t.ts.Insert(models.Timestamp{
			Timestamp: time.Unix(item.Timestamp, 0),
			Value:     item.Value,
		}); err != nil {
			t.logger.Errorf("insert data:%s", err)
			http.Error(w, httpErrors.ErrWrite.Error(), http.StatusInternalServerError)
			return
		}

		response.JsonOk(w)
	}
}

func (t *timestampHandlers) GetAverageValue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := t.ts.GetAverageValueForLast(time.Second * 60)
		if err != nil {
			t.logger.Errorf("get value: %s", err)
			http.Error(w, httpErrors.ErrGetValue.Error(), http.StatusInternalServerError)
			return
		}
		response.Json(w, map[string]interface{}{
			"result": result,
		})
	}
}
