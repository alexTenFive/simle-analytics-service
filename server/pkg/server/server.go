package server

import (
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type Server struct {
	logger *zap.SugaredLogger

	server   *http.Server
	listener net.Listener
}

func NewServer(logger *zap.SugaredLogger) *Server {
	return &Server{
		logger: logger.With("server", "ApiServer"),
	}
}

func (s *Server) Run(addr string, router *mux.Router) error {
	var err error
	s.logger.Info("starting server")

	s.server = &http.Server{
		Handler:      router,
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

func (s *Server) Shutdown() error {
	_ = s.listener.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.server.Shutdown(ctx)
}
