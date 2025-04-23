package main

import (
	"fmt"

	"github.com/nn-advith/smolurl/kvmodule"
	"github.com/nn-advith/smolurl/kvmodule/datamodel"
	logger "github.com/nn-advith/smolurl/logger"
)

func main() {

	// initialise logger
	// initialise db connection
	// start the web server
	// listen for requests

	//variables
	DBTYPE := "COUCHBASE"
	COLLECTION := "smolurl"

	//initialise logger with both stdout and logfile
	err := logger.InitLogger(true, true)
	alog := logger.GlobalLogger
	if err != nil {
		fmt.Printf("initlogger error - %v ", err)
	}
	defer logger.GlobalLogger.Close()

	//db init and defer disconnect
	dbinstance := kvmodule.InitialiseDB(DBTYPE, COLLECTION)
	defer func() {
		if err := dbinstance.Disconnect(); err != nil {
			logger.GlobalLogger.Fatal(err)
		}
	}()
	alog.Info("This is a test after global loggeer refactor")

	alog.Info("Trying insert")

	var newhash datamodel.UrlEntry = datamodel.UrlEntry{
		ID:      "SOMEHASH",
		LongURL: "https://nnadvith.netlify.app",
		Created: 45545435,
		TTL:     60000,
	}

	err = dbinstance.Insert(COLLECTION, newhash)
	if err != nil {
		logger.GlobalLogger.Error("MAIN: error during insert: ", err)
	}
}
