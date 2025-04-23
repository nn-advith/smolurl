package datamodel

type UrlEntry struct {
	ID      string `bson:"_id" json:"_id"`
	LongURL string `bson:"longurl" json:"longurl"`
	Created int64  `bson:"created" json:"created"`
	TTL     int    `bson:"ttl" json:"ttl"`
}
