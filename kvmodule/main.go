package kvmodule

import (
	"fmt"

	"github.com/nn-advith/smolurl/kvmodule/dbmodules"
	"github.com/nn-advith/smolurl/logger"
)

type DBInf interface {
	Connect() error
}

func InitialiseDB(dbtype string) {
	logger.GlobalLogger.Info("Initialising DB connection")
	var dbInstance DBInf

	switch dbtype {
	case "COUCHBASE":
		dbInstance = &dbmodules.CBConnector{
			ConnectionString: "localhost",
			Username:         "cbuser",
			Password:         "cbuser",
			Bucketname:       "Ororo",
		}
		err := dbInstance.Connect()
		if err != nil {
			logger.GlobalLogger.Fatal(fmt.Sprintf("unable to connect to couchbase: %v", err))
			return
		}
	default:
		logger.GlobalLogger.Fatal("dbtype not recognised")
	}

	fmt.Println(dbInstance)
}
