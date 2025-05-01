package middleware

import (
	"context"
	"net/http"

	"github.com/nn-advith/smolurl/kvmodule"
)

type contextKey string

const DBKey contextKey = "db"

func InjectDB(db kvmodule.DBInf) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), DBKey, db)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetDBContext(ctx context.Context) kvmodule.DBInf {
	db, ok := ctx.Value(DBKey).(kvmodule.DBInf)
	if !ok {
		return nil
	}
	return db
}
