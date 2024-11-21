package routes

// This file contains constants that are used over multiple handlers in the
// microservice

// DefaultPageSize sets the number of records returned by default from the
// database
const DefaultPageSize = 10000

// DefaultPage sets the default page which is returned. This variable is used
// with DefaultPageSize to calculate the required offset in the database
const DefaultPage = 1

// MaxPageSize sets the limit of records returned by any request
const MaxPageSize = 2500000
