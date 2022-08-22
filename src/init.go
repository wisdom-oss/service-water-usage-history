// This file contains all functions used to start the microservice. Put further prerequisites which may need to be
// initialized into this file
package main

import (
	"flag"
)

/*
Initialization Step 1 - Flag Creation

This initialization step will create a boolean flag which may trigger a healthcheck later on
*/
func init() {
	// Create a new boolean variable flag which uses an existing variable pointer for the value to be assigned
	flag.BoolVar(
		&executeHealthcheck,
		"healthcheck",
		false,
		"Run a healthcheck of the service which will check if the service can call itself and is correctly setup on the API gateway",
	)
	// Parse the created flags into their variables
	flag.Parse()
}
