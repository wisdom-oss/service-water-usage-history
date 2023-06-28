package routes

import (
	"errors"
	"net/http"
)

// BasicHandler contains just a response, that is used to show the templating
func BasicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello there"))
}

func BasicWithInternalErrorHandling(w http.ResponseWriter, r *http.Request) {
	// access the error handlers
	nativeErrorChannel := r.Context().Value("nativeErrorChannel").(chan error)
	// now publish an error to each of the wisdom errors
	nativeErrorChannel <- errors.New("native test error")
	return
}

func BasicWithWISdoMErrorHandling(w http.ResponseWriter, r *http.Request) {
	// access the error handlers
	wisdomErrorChannel := r.Context().Value("wisdomErrorChannel").(chan string)
	// now publish an error to each of the wisdom errors
	wisdomErrorChannel <- "TEMPLATE"
	return
}
