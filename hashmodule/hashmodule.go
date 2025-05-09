package hashmodule

import (
	"crypto/sha256"
	"encoding/base64"
	"strconv"
)

func GenerateHash(url string, attempts int) string {
	salted := url + strconv.Itoa(attempts)
	h := sha256.Sum256([]byte(salted))
	slug := base64.URLEncoding.EncodeToString(h[:8])
	// logger.GlobalLogger.Info(h, slug)
	return slug
}
