package controllers

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"testcase/server/internal/domain"
	"testcase/server/internal/timestamp/handlers"
)

func CreateRouter(repo domain.Repository, queue domain.Queue, log *zap.SugaredLogger) *mux.Router {
	router := mux.NewRouter()

	sendHandler := handlers.NewReceiver(queue, log)
	getAverageValuesHandler := handlers.NewGetAverageValueHandler(repo, log)

	router.HandleFunc("/send", sendHandler.Handle).Methods(http.MethodPost)
	router.HandleFunc("/avg", getAverageValuesHandler.Handle).Methods(http.MethodGet)
	return router
}
