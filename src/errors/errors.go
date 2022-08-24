// Package errors
// This package contains all errors the service may return in HTTP requests
// Those errors include unauthenticated calls and forbidden ones
// TODO: Add custom errors
package errors

import (
	"fmt"
	"net/http"

	"microservice/structs"
	"microservice/vars"
)

const UnauthorizedRequest = "UNAUTHORIZED_REQUEST"

var errorTitle = map[string]string{
	UnauthorizedRequest: "Unauthorized Request",
}

var errorDescription = map[string]string{
	UnauthorizedRequest: "The resource you tried to access requires authorization. Please check your request",
}

var httpStatus = map[string]int{
	UnauthorizedRequest: http.StatusUnauthorized,
}

func NewRequestError(errorCode string) structs.RequestError {
	return structs.RequestError{
		HttpStatus:       httpStatus[errorCode],
		HttpError:        http.StatusText(httpStatus[errorCode]),
		ErrorCode:        fmt.Sprintf("%s.%s", vars.ServiceName, errorCode),
		ErrorTitle:       errorTitle[errorCode],
		ErrorDescription: errorDescription[errorCode],
	}
}