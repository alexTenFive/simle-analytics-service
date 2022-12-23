package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"testcase/server/internal/domain"
	"testcase/server/internal/entity/models"
	"testcase/server/internal/entity/transport"
	"testcase/server/pkg/httpErrors"
	"testcase/server/pkg/response"
	"time"
)

type receiveHandler struct {
	queue domain.Queue
	log   *zap.SugaredLogger
}

func NewReceiver(queue domain.Queue, log *zap.SugaredLogger) *receiveHandler {
	return &receiveHandler{queue: queue, log: log.With("handlers", "Receiver")}
}

func (x *receiveHandler) Handle(w http.ResponseWriter, r *http.Request) {
	item := transport.SendRequest{}
	if err := item.Decode(r.Body); err != nil {
		x.log.Errorf("decode request:%s", err)
		http.Error(w, httpErrors.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}
	x.queue.Add(models.Timestamp{
		Timestamp: time.Unix(item.Timestamp, 0),
		Value:     item.Value,
	})

	response.JsonOk(w)
}
