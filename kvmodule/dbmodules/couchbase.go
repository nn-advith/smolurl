package dbmodules

import (
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/nn-advith/smolurl/logger"
)

type CBConnector struct {
	ConnectionString string
	Username         string
	Password         string
	Bucketname       string
	Cluster          *gocb.Cluster
	Bucket           *gocb.Bucket
}

func (c *CBConnector) Connect() error {
	logger.GlobalLogger.Info("Connecting to cluster and initialising bucket...")
	cluster, err := gocb.Connect("couchbase://"+c.ConnectionString, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: c.Username,
			Password: c.Password,
		},
	})
	if err != nil {
		//logger.GlobalLogger.Info(fmt.Sprintf("authentication failed:%v", err))
		return fmt.Errorf("authentication failed: %v", err)
	}
	c.Cluster = cluster
	bucket := c.Cluster.Bucket(c.Bucketname)
	c.Bucket = bucket
	err = c.Bucket.WaitUntilReady(3*time.Second, nil)
	if err != nil {
		return fmt.Errorf("bucket connection failed: %v", err)
	}
	logger.GlobalLogger.Info("Cluster connected and bucket initialised")
	return nil

}
