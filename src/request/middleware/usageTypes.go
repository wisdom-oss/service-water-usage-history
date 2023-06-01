package middleware

import (
	"context"
	"github.com/blockloop/scan/v2"
	"microservice/structs/db"
	"microservice/vars/globals"
	"microservice/vars/globals/connections"
	"net/http"
)

func AttachUsageTypes() func(handler http.Handler) http.Handler {
	logger := globals.HttpLogger
	logger.Debug().Msg("getting usage types and attaching them to the request")
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// access the request's context
			ctx := r.Context()
			// now check if the context already contains the usage types
			if ctx.Value("usageTypes") != nil {
				nextHandler.ServeHTTP(w, r)
				return
			}
			// since the usage types are not set query them from the database
			usageTypeRows, err := globals.Queries.Query(connections.DbConnection, "get-all-usage-types")
			if err != nil {
				w.Header().Add("X-Error", err.Error())
				nextHandler.ServeHTTP(w, r)
				return
			}
			var dbResults []db.DbUsageType
			err = scan.Rows(&dbResults, usageTypeRows)
			if err != nil {
				w.Header().Add("X-Error", err.Error())
				nextHandler.ServeHTTP(w, r)
				return
			}
			// now append the usage types to the request's context
			ctx = context.WithValue(ctx, "usageTypes", dbResults)
			r = r.WithContext(ctx)
			// now serve the request to the next handler
			nextHandler.ServeHTTP(w, r)
		})
	}
}
