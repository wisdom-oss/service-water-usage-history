// This file contains all globally used variables and their default values
package vars

import (
	"os"

	"microservice/structs"
)

// TODO: Change me
const ServiceName = "change-me"

var (
	// The host on which the API gateway runs and this service shall be registered on
	ApiGatewayHost = os.Getenv("CONFIG_API_GATEWAY_HOST")
	// The administration port of the api gateway on the host
	ApiGatewayAdminPort = os.Getenv("CONFIG_API_GATEWAY_ADMIN_PORT")
	// The path on which the service shall be reachable
	ApiGatewayServicePath = os.Getenv("CONFIG_API_GATEWAY_SERVICE_PATH")
	// The http port on which the service will listen for new requests
	HttpListenPort = os.Getenv("CONFIG_HTTP_LISTEN_PORT")
	// Indicator if a health check shall be executed instead of the main() function
	ExecuteHealthcheck = false
	// The scope the service was configured with
	Scope *structs.ScopeInformation
)
