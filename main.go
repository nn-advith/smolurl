package main

import (
	"fmt"

	"github.com/nn-advith/smolurl/kvmodule"
	logger "github.com/nn-advith/smolurl/logger"
)

func main() {

	// initialise logger
	// initialise db connection
	// start the web server
	// listen for requests

	err := logger.InitLogger(true, true)
	alog := logger.GlobalLogger
	if err != nil {
		fmt.Printf("initlogger error - %v ", err)
	}
	defer logger.GlobalLogger.Close()

	kvmodule.InitialiseDB("COUCHBASE")

	alog.Info("This is a test after global loggeer refactor")
}
