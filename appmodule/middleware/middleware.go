package middleware

import (
	"context"
	"net/http"

	"github.com/nn-advith/smolurl/kvmodule"
)

type contextKey string

const DBKey contextKey = "db"

type DBMiddleware struct {
	handler http.Handler
	db      kvmodule.DBInf
}

func (dm *DBMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), DBKey, dm.db)
	dm.handler.ServeHTTP(w, r.WithContext(ctx))
}

func NewDBMiddleware(handlerToWrap http.Handler, db kvmodule.DBInf) *DBMiddleware {
	return &DBMiddleware{handler: handlerToWrap, db: db}
}

func GetDBContext(ctx context.Context) kvmodule.DBInf {
	db, ok := ctx.Value(DBKey).(kvmodule.DBInf)
	if !ok {
		return nil
	}
	return db
}
