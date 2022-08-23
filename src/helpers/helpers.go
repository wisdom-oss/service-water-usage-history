package helpers

import (
	"fmt"
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
