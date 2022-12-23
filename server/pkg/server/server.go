package server

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type Server struct {
	db     *sqlx.DB
	logger *zap.SugaredLogger

	server   *http.Server
	listener net.Listener
}

func NewServer(db *sqlx.DB, logger *zap.SugaredLogger) *Server {
	return &Server{
		db:     db,
		logger: logger.With("server", "ApiServer"),
	}
}

func (s *Server) Run(addr string, router *mux.Router) error {
	var err error
	s.logger.Info("starting server")

	s.server = &http.Server{
		Handler: router,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	go func() {
		if err = s.server.Serve(s.listener); err != nil {
			s.logger.Errorf("Serve: %s", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown() {
	s.listener.Close()
	s.server.Close()
}
