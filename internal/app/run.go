package app

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"testcase_v2/internal/controllers/http"
	"testcase_v2/internal/usecase"
	"testcase_v2/internal/usecase/queue"
	"testcase_v2/internal/usecase/repository"
	"testcase_v2/pkg/db/postgres"
	"testcase_v2/pkg/logger"
	"testcase_v2/pkg/server"
)

func Run() {
	log := logger.Init(true)
	pg, err := postgres.NewPgDB(
		viper.GetString("postgres.host"),
		viper.GetString("postgres.port"),
		viper.GetString("postgres.user"),
		viper.GetString("postgres.dbname"),
		viper.GetString("postgres.password"),
	)
	if err != nil {
		log.Errorf("init database: %s", err)
		os.Exit(1)
	}

	tsrepo := repository.NewDatabaseRepository(pg, log)
	exitCh := make(chan chan struct{})
	tsqueue := queue.NewRamQueue(tsrepo, exitCh, log)
	router := http.CreateRouter(usecase.NewTimestampUseCase(tsrepo, tsqueue, log), log)

	srv := server.NewServer(log)
	if err = srv.Run(fmt.Sprintf(":%s", viper.GetString("http.port")), router); err != nil {
		log.Errorf("server.Run(): %s", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	if err = srv.Shutdown(); err != nil {
		log.Errorf("server.Shutdown(): %s", err)
	}
	// wait for processing tsqueue
	confirmExitCh := make(chan struct{})
	exitCh <- confirmExitCh
	<-confirmExitCh

	log.Info("server turned off")
}
