package appmodule

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nn-advith/smolurl/appmodule/middleware"
	"github.com/nn-advith/smolurl/appmodule/server"
	"github.com/nn-advith/smolurl/kvmodule"
	"github.com/nn-advith/smolurl/logger"
)

func generateRoutes(db kvmodule.DBInf) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		dbinstance := middleware.GetDBContext(req.Context())
		data, err := dbinstance.Read("smolurl", "SOMEHASH")
		if err != nil {
			logger.GlobalLogger.Error("MAIN: error during Read: ", err)
		}

		fmt.Fprintln(w, data)
	})

	newMux := middleware.NewDBMiddleware(mux, db)
	return newMux

}

func ConfigureAppModule(dbinstance kvmodule.DBInf) {

	cfg := server.Config{
		Address:      ":4000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
	m := generateRoutes(dbinstance)
	app, err := server.NewServer(cfg, dbinstance, m)
	if err != nil {
		logger.GlobalLogger.Fatal("Error creating new server")
	}

	go func() {
		if err := app.Start(); err != nil && err != http.ErrServerClosed {
			fmt.Println("start error: ", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		fmt.Println("Shutdown error:", err)
	}

	logger.GlobalLogger.Info("Server gracefully stopped")
}
