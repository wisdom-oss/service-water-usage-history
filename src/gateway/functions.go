package gateway

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"microservice/helpers"
	"net/http"
	"net/url"
)

var gatewayAPIUrl = ""
var gatewayUpstreamName = ""
var gatewayServiceName = ""
var connectionsPrepared = false

var logger = log.WithFields(log.Fields{
	"localIP":             helpers.GetLocalIP(),
	"gatewayAPIUrl":       gatewayAPIUrl,
	"gatewayUpstreamName": gatewayUpstreamName,
	"gatewayServiceName":  gatewayServiceName,
})

// Upstream
// Information about the upstream used by this service
var Upstream *UpstreamConfiguration

/*
PrepareGatewayConnections

Call this function once before calling the api gateway to prepare the gatewayAPIUrl. T
*/
func PrepareGatewayConnections(serviceName string, host string, adminAPIPort string) {
	gatewayAPIUrl = fmt.Sprintf("http://%s:%s", host, adminAPIPort)
	logger.Info("Set the gatewayAPIUrl")
	gatewayUpstreamName = fmt.Sprintf("upstream_%s", serviceName)
	logger.Info("Successfully set the gateway upstream name")
	gatewayServiceName = fmt.Sprintf("service_%s", serviceName)
	logger.Info("Successfully set the gateway service name")
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
		logger.WithField("function", "CreateUpstream").WithError(err).Error("An error occurred while sending the request to create the new upstream")
		return false
	}
	if response.StatusCode != 201 {
		logger.WithField("function", "CreateUpstream").Error("The gateway did not report the correct response code. The upstream may not have been created")
		return false
	}
	logger.WithField("function", "CreateUpstream").Info("The upstream was created in the gateway. Storing information about the upstream.")
	parseError := json.NewDecoder(response.Body).Decode(&Upstream)
	if parseError != nil {
		logger.WithField("function", "CreateUpstream").Warning("Unable to parse the content of the response. There will be no information available about the upstream")
	}
	return true
}
