package appmodule

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nn-advith/smolurl/appmodule/middleware"
	"github.com/nn-advith/smolurl/appmodule/server"
	"github.com/nn-advith/smolurl/hashmodule"
	"github.com/nn-advith/smolurl/kvmodule"
	"github.com/nn-advith/smolurl/kvmodule/datamodel"
	"github.com/nn-advith/smolurl/logger"
)

type Rbody struct {
	Url string `json:"url"`
}

var COLLECTION string

func generateRoutes(db kvmodule.DBInf) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("standard"))
	})

	mux.HandleFunc("POST /hash/", func(w http.ResponseWriter, r *http.Request) {
		dbinstance := middleware.GetDBContext(r.Context())
		var rbody Rbody
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&rbody)
		if err != nil {
			logger.GlobalLogger.Error("unable to decode request body; dropping")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//pass the url to hash module and generate a hash
		//query the database for hash collisions; if yes, retrigger hash generation; else add to db with related fields
		//respond with hash url as plain text
		attempts := 1
		slug := hashmodule.GenerateHash(rbody.Url, attempts)
		for {
			urlentry, err := dbinstance.Read(COLLECTION, slug)
			if err != nil {
				logger.GlobalLogger.Error("smallgen error: database read: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if urlentry == nil {
				break
			}

			urlasserted := urlentry.(datamodel.UrlEntry)
			if urlasserted.LongURL == rbody.Url {
				//exists
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("http://localhost:4000/%v\n", slug)))
				return
			} else {
				attempts++
				slug = hashmodule.GenerateHash(rbody.Url, attempts)
			}
		}
		//slug is ready and no collisions
		//add to db
		var newUrlEntry datamodel.UrlEntry
		newUrlEntry.Created = time.Now().Unix()
		newUrlEntry.ID = slug
		newUrlEntry.TTL = 60
		newUrlEntry.LongURL = rbody.Url

		err = dbinstance.Insert(COLLECTION, newUrlEntry)
		if err != nil {
			logger.GlobalLogger.Error("smolgen error: database insert", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(fmt.Sprintf("http://localhost:4000/%v\n", slug)))

	})

	mux.HandleFunc("GET /{somehash}", func(w http.ResponseWriter, r *http.Request) {
		dbinstance := middleware.GetDBContext(r.Context())
		hashval := r.PathValue("somehash")
		urlenty, err := dbinstance.Read(COLLECTION, hashval)
		if err != nil {
			logger.GlobalLogger.Error("queryhash error: ", err)
		}
		if urlenty == nil {
			//noentry
			w.WriteHeader(http.StatusNotFound)
		} else {
			urlasserted := urlenty.(datamodel.UrlEntry)
			if time.Now().Unix() > (urlasserted.Created + int64(urlasserted.TTL)) {
				//delete entry
				err := dbinstance.Delete(COLLECTION, urlasserted.ID)
				if err != nil {
					logger.GlobalLogger.Error("queryhash error: ", err)
				}
				w.WriteHeader(http.StatusNotFound)
				return
			}
			http.Redirect(w, r, urlasserted.LongURL, http.StatusMovedPermanently)
		}

	})

	newMux := middleware.NewDBMiddleware(mux, db)
	return newMux

}

func ConfigureAppModule(dbinstance kvmodule.DBInf, collection string) {
	COLLECTION = collection
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
