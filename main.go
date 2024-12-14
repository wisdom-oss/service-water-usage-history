package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	healthcheckServer "github.com/wisdom-oss/go-healthcheck/server"

	"microservice/internal"
	"microservice/internal/config"
	"microservice/internal/db"
	"microservice/routes"

	"github.com/wisdom-oss/common-go/v3/middleware/gin/jwt"
)

// the main function bootstraps the http server and handlers used for this
// microservice
func main() {
	// create a new logger for the main function
	l := log.Logger
	l.Info().Msgf("configuring %s service", internal.ServiceName)

	// create the healthcheck server
	hcServer := healthcheckServer.HealthcheckServer{}
	hcServer.InitWithFunc(func() error {
		// test if the database is reachable
		return db.Pool.Ping(context.Background())
	})
	err := hcServer.Start()
	if err != nil {
		l.Fatal().Err(err).Msg("unable to start healthcheck server")
	}
	go hcServer.Run()

	// prepare some scope requirers to make the route definition easiser
	scopeRequirer := jwt.ScopeRequirer{}
	scopeRequirer.Configure(internal.ServiceName)

	r := config.PrepareRouter()
	r.GET("/", scopeRequirer.RequireRead, routes.PagedUsages)

	// create http server
	server := &http.Server{
		Addr:    config.ListenAddress,
		Handler: r,
	}

	l.Info().Msg("starting http server")

	// Start the server and log errors that happen while running it
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			l.Fatal().Err(err).Msg("An error occurred while starting the http server")
		}
	}()

	// Set up some the signal handling to allow the server to shut down gracefully
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	// Block further code execution until the shutdown signal was received
	l.Info().Msg("server ready to accept connections")
	<-shutdownSignal

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		l.Fatal().Err(err).Msg("An error occurred while shutting down http server")
	}

}
