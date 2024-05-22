//go:build !docker

package config

import (
	"net/http"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	errorMiddleware "github.com/wisdom-oss/microservice-middlewares/v5/error"
)

// This file contains default paths that are used inside the service to load
// resource files.
// The go toolchain automatically uses this file if no build tags have been
// set.
// This should be the case in local development and testing.
// To achieve the behavior in docker containers, supply the build tag "docker"
// and the Docker defaults are used

// DefaultMiddlewares contains the middlewares used per default in this service
var DefaultMiddlewares = []func(next http.Handler) http.Handler{
	chiMiddleware.RequestID,
	chiMiddleware.RealIP,
	errorMiddleware.Handler,
}
