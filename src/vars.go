// This file contains all globally used variables and their default values
package main

import "os"

var (
	// The host on which the API gateway runs and this service shall be registered on
	apiGatewayHost = os.Getenv("CONFIG_API_GATEWAY_HOST")
	// The administration port of the api gateway on the host
	apiGatewayAdminPort = os.Getenv("CONFIG_API_GATEWAY_ADMIN_PORT")
	// The path on which the service shall be reachable
	apiGatewayServicePath = os.Getenv("CONFIG_API_GATEWAY_SERVICE_PATH")
	// The http port on which the service will listen for new requests
	httpListenPort = os.Getenv("CONFIG_HTTP_LISTEN_PORT")
	// Indicator if a health check shall be executed instead of the main() function
	executeHealthcheck = false
)
