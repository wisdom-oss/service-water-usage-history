//go:build docker

package config

import (
	"net/http"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// This file contains default paths that are used inside the service to load
// resource files.
// The go toolchain automatically uses this file if the build tag "docker" has
// been set.
// This should be the case in docker images and containers only.
// To achieve the behavior for local development and testing, please remove the
// "docker" build tag

var DefaultMiddlewares = []func(next http.Handler) http.Handler{
	chiMiddleware.RequestID,
	chiMiddleware.RealIP,
	errorMiddleware.Handler,
	securityMiddleware.ValidateServiceJWT,
}
