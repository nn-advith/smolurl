package kvmodule

import (
	"fmt"

	"github.com/nn-advith/smolurl/kvmodule/dbmodules"
	"github.com/nn-advith/smolurl/logger"
)

type DBInf interface {
	Connect(collection string) error
	Disconnect() error
	Insert(collection string, data any) error
}

func InitialiseDB(dbtype string, collection string) DBInf {
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
		err := dbInstance.Connect(collection)
		if err != nil {
			logger.GlobalLogger.Fatal(fmt.Sprintf("unable to connect to couchbase: %v", err))
			return nil
		}
	default:
		logger.GlobalLogger.Fatal("dbtype not recognised")
	}

	return dbInstance
}
