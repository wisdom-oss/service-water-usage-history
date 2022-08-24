package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
	"microservice/helpers"
)

var gatewayAPIUrl = ""
var gatewayUpstreamName = ""
var gatewayServiceName = ""
var connectionsPrepared = false
var httpListenPort = ""

var logger = log.WithFields(log.Fields{
	"localIP":             helpers.GetLocalIP(),
	"gatewayAPIUrl":       gatewayAPIUrl,
	"gatewayUpstreamName": gatewayUpstreamName,
	"gatewayServiceName":  gatewayServiceName,
})

// Upstream
// Information about the upstream used by this service
var Upstream *UpstreamConfiguration
var ServiceEntry *ServiceConfiguration

/*
PrepareGatewayConnections

Call this function once before calling the api gateway to prepare the gatewayAPIUrl. T
*/
func PrepareGatewayConnections(serviceName string, host string, adminAPIPort string, listenPort string) {
	gatewayAPIUrl = fmt.Sprintf("http://%s:%s", host, adminAPIPort)
	logger.Info("Set the gatewayAPIUrl")
	gatewayUpstreamName = fmt.Sprintf("upstream_%s", serviceName)
	logger.Info("Successfully set the gateway upstream name")
	gatewayServiceName = fmt.Sprintf("service_%s", serviceName)
	logger.Info("Successfully set the gateway service name")
	httpListenPort = listenPort
	connectionsPrepared = true
}

/*
IsUpstreamSetUp

Check if the upstream that will be used by the service is already configured.

If it is configured the information about the upstream will be stored.
*/
func IsUpstreamSetUp() bool {
	if !connectionsPrepared {
		logger.WithField("function", "IsUpstreamSetUp").Warning("The gateway connections have not been prepared before calling this method")
		return false
	}
	// Request information from the api gateway about the upstream this service will use
	response, err := http.Get(gatewayAPIUrl + "/upstreams/" + gatewayUpstreamName)
	if err != nil {
		logger.WithField("function", "IsUpstreamSetUp").WithError(err).Error("An error occurred while requesting information about the upstream used by the service")
		return false
	}
	switch response.StatusCode {
	case 200:
		logger.WithField("function", "IsUpstreamSetUp").Info("The upstream used by the service was found in the api gateway")
		parseError := json.NewDecoder(response.Body).Decode(&Upstream)
		if parseError != nil {
			logger.WithField("function", "IsUpstreamSetUp").WithError(err).Error("Unable to parse the content of the upstream information request")
		}
		return true
	case 404:
		logger.WithField("function", "IsUpstreamSetup").Warning("The upstream used by the service was not found in the api gateway. Consider creating it to allow the service to work.")
		return false
	default:
		return false
	}
}

/*
CreateUpstream

Create a new upstream in the api gateway for this service and its instances
*/
func CreateUpstream() bool {
	if !connectionsPrepared {
		logger.WithField("function", "CreateUpstream").Warning("The gateway connections have not been prepared before calling this method")
	}
	// Build the request content
	requestBody := url.Values{}
	requestBody.Set("name", gatewayUpstreamName)
	// Post the request body to the gateway
	response, err := http.PostForm(gatewayAPIUrl+"/upstreams", requestBody)
	if err != nil {
		logger.
			WithField("function", "CreateUpstream").
			WithError(err).
			Error("An error occurred while sending the request to create the new upstream")
		return false
	}
	if response.StatusCode != 201 {
		logger.
			WithField("function", "CreateUpstream").
			Error("The gateway did not report the correct response code. The upstream may not have been created")
		return false
	}
	logger.
		WithField("function", "CreateUpstream").
		Info("The upstream was created in the gateway. Storing information about the upstream.")
	parseError := json.NewDecoder(response.Body).Decode(&Upstream)
	if parseError != nil {
		logger.
			WithField("function", "CreateUpstream").
			Warning("Unable to parse the content of the response. " +
				"There will be no information available about the upstream")
	}
	return true
}

/*
IsIPAddressInUpstreamTargets

Check if this instance of the service is in the target list of the upstream for this service
*/
func IsIPAddressInUpstreamTargets() bool {
	logger := logger.WithField("function", "IsIPAddressInUpstreamTargets")
	if !connectionsPrepared {
		logger.
			Warning("The gateway connections have not been prepared before calling this method")
	}
	// Get the local IP address of the container
	localIPAddress := helpers.GetLocalIP()
	targetAddress := localIPAddress + ":" + httpListenPort
	logger.
		WithField("targetAddress", targetAddress).
		Info("Checking if this service instance is registered in the upstream for this service")
	// Get a list of the targets configured for the upstream
	response, err := http.Get(gatewayAPIUrl + "/upstreams/" + Upstream.Id + "/targets")
	if err != nil {
		logger.
			WithError(err).
			Error("An error occurred while requesting the list of targets for the current upstream")
		return false
	}
	targetList := new(TargetList)
	parserError := json.NewDecoder(response.Body).Decode(&targetList)
	if parserError != nil {
		logger.
			WithError(parserError).
			Error("Unable to parse the target list from the response sent by the gateway")
		return false
	}
	for _, target := range targetList.Targets {
		if target.Address == targetAddress {
			logger.Info("Found the target in the current upstream")
			return true
		}
	}
	logger.Warning("The target is not in the current upstream")
	return false
}

func AddServiceToUpstreamTargets() bool {
	logger := logger.WithFields(log.Fields{
		"function":   "AddServiceToUpstreamTargets",
		"upstreamId": Upstream.Id,
	})
	if !connectionsPrepared {
		logger.
			Warning("The gateway connections have not been prepared before calling this method")
	}
	// Get the local IP address of the container
	localIPAddress := helpers.GetLocalIP()
	targetAddress := localIPAddress + ":" + httpListenPort
	logger.
		WithField("targetAddress", targetAddress).
		Info("Registering the service instance as target in the upstream for this service")
	// Create the request body
	requestBody := url.Values{}
	requestBody.Set("target", targetAddress)
	// Send the request body to the gateway
	response, err := http.PostForm(gatewayAPIUrl+"/upstreams/"+Upstream.Id+"/targets", requestBody)
	if err != nil {
		logger.
			WithError(err).
			Error("An error occurred while sending the request to the api gateway")
		return false
	}
	// Now check if the response code is 201 indicating that the target has been created
	if response.StatusCode != 201 {
		logger.
			Warning("The gateway did not respond with a 201 Created to the request. " +
				"The target may not have been created. Use the function 'gateway.IsIPAddressInUpstreamTargets(" +
				")' to test again")
		return false
	}
	logger.Info("The target was successfully created in the upstream")
	return true
}

/*
IsServiceSetUp

Check if a service was created in the gateway.
A service is an object in the gateway which may be used to create new routes
*/
func IsServiceSetUp() bool {
	logger := logger.WithFields(log.Fields{
		"function": "IsServiceSetUp",
	})
	if !connectionsPrepared {
		logger.
			Warning("The gateway connections have not been prepared before calling this method")
	}
	logger.Info("Checking is a service entry exists on the gateway")
	response, err := http.Get(gatewayAPIUrl + "/services/" + gatewayServiceName)
	if err != nil {
		logger.WithError(err).Error("An error occurred while sending the request to the api gateway")
	}
	if response.StatusCode != 200 {
		logger.Warning("There is no service entry for this service in the gateway")
		return false
	}
	decodeErr := json.NewDecoder(response.Body).Decode(&ServiceEntry)
	if decodeErr != nil {
		logger.WithError(decodeErr).Warning("Unable to parse the service information sent by the gateway")
		return true
	}
	logger.Info("Found a service entry in the gateway")
	return true
}
