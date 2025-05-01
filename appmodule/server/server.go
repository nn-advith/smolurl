package server

import (
	"context"
	"net/http"
	"time"

	"github.com/nn-advith/smolurl/kvmodule"
	"github.com/nn-advith/smolurl/logger"
)

type Config struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type AppServer struct {
	Server     http.Server
	DBInstance *kvmodule.DBInf
}

func NewServer(cfg Config, dbInstance kvmodule.DBInf, handler http.Handler) (*AppServer, error) {
	return &AppServer{
		Server: http.Server{
			Addr:         cfg.Address,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		DBInstance: &dbInstance,
	}, nil
}

func (a *AppServer) Start() error {
	logger.GlobalLogger.Info("staring server on address: ", a.Server.Addr)
	return a.Server.ListenAndServe()
}

func (a *AppServer) Stop(ctx context.Context) error {
	return a.Server.Shutdown(ctx)
}
