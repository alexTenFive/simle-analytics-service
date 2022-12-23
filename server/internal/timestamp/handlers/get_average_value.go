package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"testcase/server/internal/domain"
	"testcase/server/internal/timestamp/usecase"
	"testcase/server/pkg/httpErrors"
	"testcase/server/pkg/response"
	"time"
)

type getAverageValueHandler struct {
	uc     domain.GetAverageValueUseCase
	logger *zap.SugaredLogger
}

func (x *getAverageValueHandler) Handle(w http.ResponseWriter, r *http.Request) {
	result, err := x.uc.Do(time.Second * 60)
	if err != nil {
		x.logger.Errorf("get value: %s", err)
		http.Error(w, httpErrors.ErrGetValue.Error(), http.StatusInternalServerError)
		return
	}
	response.Json(w, map[string]interface{}{
		"result": result,
	})
}

func NewGetAverageValueHandler(repo domain.Repository, log *zap.SugaredLogger) *getAverageValueHandler {
	return &getAverageValueHandler{
		uc:     usecase.NewGetAverageValueUseCase(repo),
		logger: log.With("handler", "GetAverageValue"),
	}
}
