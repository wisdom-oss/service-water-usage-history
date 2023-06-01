package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	_ "github.com/wisdom-oss/microservice-middlewares"
	customMiddleware "microservice/request/middleware"
	"microservice/request/routes"
	"microservice/vars/globals"
	"net/http"
	"os"
	"os/signal"
	"time"
)

/*
This function is used to set up the http server for the microservice
*/
func main() {

	// Set up the routing of the different functions
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(httplog.Handler(globals.HttpLogger))
	router.Use(middleware.Heartbeat("/healthcheck"))
	router.Use(middleware.Compress(5))
	//	router.Use(wisdomMiddleware.Authorization([]string{"/healthcheck"}, globals.ScopeConfiguration.ScopeValue))
	router.Use(customMiddleware.AttachUsageTypes())
	router.HandleFunc("/all", routes.GetAllUsages)
	router.HandleFunc("/by-consumer/{consumerId}", routes.GetConsumerUsages)

	// Configure the HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", globals.Environment["LISTEN_PORT"]),
		WriteTimeout: 6 * time.Hour,
		ReadTimeout:  6 * time.Hour,
		IdleTimeout:  6 * time.Hour,
		Handler:      router,
	}

	// Start the server and log errors that happen while running it
	go func() {
		if err := server.ListenAndServe(); err != nil {
			globals.HttpLogger.Fatal().Err(err).Msg("An error occurred while starting the http server")
		}
	}()

	globals.HttpLogger.Info().Msg("server ready to accept connections")

	// Set up the signal handling to allow the server to shut down gracefully
	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, os.Interrupt)

	// Block further code execution until the shutdown signal was received
	<-cancelSignal

}
