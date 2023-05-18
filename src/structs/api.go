package structs

import "github.com/paulmach/go.geojson"

// IncomingUsageRecord reflects the expected input when creating a new water
// usage record using this microservice
type IncomingUsageRecord struct {
	// Date contains a UNIX timestamp without microseconds indicating
	// into which year the usage may be counted.
	Date int64 `json:"date"`

	// UsageType contains the string external identifier of a usage type which
	// indicates the type of usage for this record
	UsageType *string `json:"usageType"`

	// Consumer contains the uuid identifier of a consumer which is
	// assigned to the usage record
	Consumer *string `json:"consumer"`

	// Amount contains the used water amount in cubic meters
	Amount float64 `json:"amount"`
}

// UsageRecord extends the incoming usage record by fields that were either
// generated in the database or added by other processes
type UsageRecord struct {
	IncomingUsageRecord

	// RecordID contains the internal ID of the usage record which may
	// be used to remove a usage record from the history
	RecordID int64 `json:"recordID"`

	// RecordedAt contains a timestamp indicating the time the record was
	// written into the database. To get the time for which the record was made
	// use the Date attribute of this object
	RecordedAt int64 `json:"recordedAt"`
}

// IncomingUsageType reflects the expected input to the api when creating a new
// usage type in the database
type IncomingUsageType struct {
	// Title contains a short title for the new usage type and is required from
	// the database
	Title string `json:"title"`

	// Description contains an optional description for the usage type
	Description *string `json:"description"`

	// ExternalIdentifier contains the required external identifier for this
	// usage type that is used to identify this usage type in queries to the api
	ExternalIdentifier string `json:"externalIdentifier"`
}

// UsageType extends the IncomingUsageType by field that are created in the
// database. This struct should be used to return information about a usage type
// in responses to api requests
type UsageType struct {
	IncomingUsageRecord

	// ID contains the generated UUID from the database which allows the
	// identification of this usage type when modifying or deleting it
	ID string `json:"id"`
}

// IncomingConsumer reflects the expected input into the api when creating a new
// consumer via this microservice
type IncomingConsumer struct {
	// Name contains the clear name of the consumer which shall be created
	Name string `json:"name"`

	// Coordinates contains an array consisting of two floats which indicate the
	// latitude and longitude of the consumer's location. The first element
	// is the latitude and the second element is the longitude
	Coordinates [2]float64 `json:"coordinates"`

	// UsageType contains an optional externalIdentifier of a usage type which
	// can be associated to a consumer to automatically setting the usage type
	// when recording a new usage for a consumer
	UsageType *string `json:"usageType"`

	// AdditionalProperties contains an optional key/value map allowing to add
	// additional properties to a consumer
	AdditionalProperties map[string]interface{} `json:"additionalProperties"`
}

// Consumer is a struct that should be used to return information about a
// consumer which was stored in the database.
type Consumer struct {
	// ID contains the database-generated UUID for this consumer
	ID string `json:"ID"`

	// Name contains the name of the consumer
	Name string `json:"name"`

	// Location contains the GeoJSON representation of the consumers location
	Location geojson.Geometry `json:"location"`

	// UsageType contains the UUID of the usage type that was assigned to the
	// consumer
	UsageType *string `json:"usageType"`

	// AdditionalProperties contains an optional key/value map allowing to add
	// additional properties to a consumer
	AdditionalProperties map[string]interface{} `json:"additionalProperties"`
}
