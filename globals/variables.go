package globals

import (
	_ "embed"

	"github.com/qustavo/dotsql"
)

// This file contains globally shared variables (e.g., service name, sql queries)

// ServiceName contains the service's identifying name.
// It is automatically loaded from the service-name file in the repository.
//
//go:embed service-name
var ServiceName string

// SqlQueries contains the prepared sql queries from the resources folder
var SqlQueries *dotsql.DotSql

// Environment contains a mapping between the environment variables and the values
// they were set to. However, this variable only contains the configured environment
// variables
var Environment map[string]string = make(map[string]string)
