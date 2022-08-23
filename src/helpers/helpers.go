package helpers

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

/*
PingHost

Use this function to ping a host on a port by connecting to it via a tcp connection.

Parameters:
	- host: The host which is the target of the connection
	- port: The port which is the target of the connection
	- timeout: The connection timeout in seconds

The function returns true if the connection was successful. Else it will return false
*/
func PingHost(host string, port string, timeout int) bool {
	// Build the connection timeout
	connectionTimeout := time.Duration(timeout) * time.Second
	// Build the tcp connection target string
	connectionTarget := fmt.Sprintf("%s:%s", host, port)
	// Open a tcp connection with the parameters
	_, connectionError := net.DialTimeout("tcp", connectionTarget, connectionTimeout)
	if connectionError != nil {
		return false
	} else {
		return true
	}

}

/*
GetLocalIP

Get the local ip address as a string.
*/
func GetLocalIP() string {
	interfaceAddresses, err := net.InterfaceAddrs()
	if err != nil {
		log.WithError(err).Error("Unable to access the network interface addresses. Please check your permissions.")
		return ""
	}
	for _, address := range interfaceAddresses {
		if address, ok := address.(*net.IPNet); ok && !address.IP.IsLoopback() {
			if address.IP.To4() != nil {
				return address.IP.String()
			}
		}
	}
	return ""
}