package main

import (
	"os"
	"os/signal"
	"syscall"
	"testcase/server/internal/controllers"
	"testcase/server/internal/data/queue"
	timestampRepo "testcase/server/internal/data/repository"
	"testcase/server/pkg/db/postgres"
	"testcase/server/pkg/logger"
	"testcase/server/pkg/server"
)

const apiBase = "127.0.0.1:8888"

func main() {
	log := logger.Init(true)
	pgdb, err := postgres.NewPgDB("localhost", "54320", "postgres", "timestamps", "secret")
	if err != nil {
		log.Errorf("init database: %s", err)
		os.Exit(1)
	}
	exitCh := make(chan chan struct{})

	repo := timestampRepo.NewDatabaseRepository(pgdb, log)
	router := controllers.CreateRouter(repo, queue.NewRamQueue(repo, exitCh, log), log)

	srv := server.NewServer(pgdb, log)
	if err = srv.Run(apiBase, router); err != nil {
		log.Errorf("server.Run(): %s", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	srv.Shutdown()
	// wait for processing rest of queries
	confirmExitCh := make(chan struct{})
	exitCh <- confirmExitCh
	<-confirmExitCh

	log.Info("server exited properly")
}
