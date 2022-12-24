package http

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"testcase_v2/internal/usecase"
)

func CreateRouter(ts usecase.TimestampUseCase, log *zap.SugaredLogger) *mux.Router {
	router := mux.NewRouter()

	tsHandlers := newTimestampHandlers(ts, log)
	router.HandleFunc("/send", tsHandlers.Receive()).Methods(http.MethodPost)
	router.HandleFunc("/avg", tsHandlers.GetAverageValue()).Methods(http.MethodGet)

	return router
}
