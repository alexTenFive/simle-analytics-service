package app

import (
	"os"
	"os/signal"
	"syscall"
	"testcase_v2/server/internal/controllers/http"
	"testcase_v2/server/internal/usecase"
	"testcase_v2/server/internal/usecase/queue"
	"testcase_v2/server/internal/usecase/repository"
	"testcase_v2/server/pkg/db/postgres"
	"testcase_v2/server/pkg/logger"
	"testcase_v2/server/pkg/server"
)

const apiBase = "127.0.0.1:8888"

func Run() {
	log := logger.Init(true)
	pg, err := postgres.NewPgDB("localhost", "54320", "postgres", "timestamps", "secret")
	if err != nil {
		log.Errorf("init database: %s", err)
		os.Exit(1)
	}
	exitCh := make(chan chan struct{})

	tsrepo := repository.NewDatabaseRepository(pg, log)
	tsqueue := queue.NewRamQueue(tsrepo, exitCh, log)

	router := http.CreateRouter(usecase.NewTimestampUseCase(tsrepo, tsqueue, log), log)

	srv := server.NewServer(pg, log)
	if err = srv.Run(apiBase, router); err != nil {
		log.Errorf("server.Run(): %s", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	srv.Shutdown()
	// wait for processing tsqueue
	confirmExitCh := make(chan struct{})
	exitCh <- confirmExitCh
	<-confirmExitCh

	log.Info("server exited properly")
}
