package dbmodules

import (
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/nn-advith/smolurl/kvmodule/datamodel"
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

func (c *CBConnector) Connect(collection string) error {
	logger.GlobalLogger.Info("CB: Connecting to cluster and initialising bucket...")
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

	// check for existance of collection

	cMgr := c.Bucket.CollectionsV2()
	scopes, err := cMgr.GetAllScopes(nil)

	if err != nil {
		return fmt.Errorf("unable to get all scopes for bucket: %v", err)
	}

	exists := false
	for _, s := range scopes {
		if s.Name == "_default" {
			for _, c := range s.Collections {
				if c.Name == collection {
					exists = true
					break
				}
			}
		}
	}
	if !exists {
		//create collection under default scope
		// err := cMgr.CreateCollection("_default", collection, nil, nil)
		// if err != nil {
		// 	return fmt.Errorf("unable to create collection; pls verify couchbase : %v", err)
		// }
		// logger.GlobalLogger.Info("CB: Couchbase collection ", collection, " created")

		return fmt.Errorf("required collection is not present, please contact admin to create collection")
	}
	logger.GlobalLogger.Info("CB: Collection found.")
	logger.GlobalLogger.Info("CB: Cluster connected and bucket initialised")
	return nil

}

func (c *CBConnector) Disconnect() error {
	if err := c.Cluster.Close(nil); err != nil {
		return fmt.Errorf("error during disconnect: %v", err)
	} else {
		logger.GlobalLogger.Info("CB: disconnected from couchbase")
	}
	return nil
}

func (c *CBConnector) Insert(collection string, data any) error {
	logger.GlobalLogger.Info("CB: Inserting document")
	col := c.Bucket.DefaultScope().Collection(collection)
	if e, ok := data.(datamodel.UrlEntry); ok { //type-assertion - move to generics
		_, err := col.Insert(e.ID, data, nil)
		if err != nil {
			return fmt.Errorf("error during insert: %v", err)
		}
	} else {
		return fmt.Errorf("data is not of type UrlEntry")
	}
	return nil
}

func (c *CBConnector) Update(collection string, data any) error {
	logger.GlobalLogger.Info("CB: Updating document")
	col := c.Bucket.DefaultScope().Collection(collection)
	if e, ok := data.(datamodel.UrlEntry); ok {
		_, err := col.Replace(e.ID, data, nil)
		if err != nil {
			return fmt.Errorf("error during update: %v", err)
		}
	} else {
		return fmt.Errorf("data is not of type UrlEntry")
	}
	return nil
}

func (c *CBConnector) Read(collection string, id string) (any, error) {
	logger.GlobalLogger.Info("CB: Retreiving document with ID ", id)
	col := c.Bucket.DefaultScope().Collection(collection)
	doc, err := col.Get(id, nil)
	if err != nil {
		return nil, fmt.Errorf("error during retreival: %v", err)
	}
	var urlentry datamodel.UrlEntry
	if err := doc.Content(&urlentry); err != nil {
		return nil, fmt.Errorf("unable to decode: %v", err)
	}
	return urlentry, nil
}
