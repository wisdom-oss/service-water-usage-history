package structs

import (
	"github.com/jackc/pgtype"
	geojson "github.com/paulmach/go.geojson"
)

// DbWaterUsageRecord reflects a single water usage record stored in the
// WISdoM database. It may be converted to a json-encodable object using
// the ToUsageRecord() function on a given instance.
type DbWaterUsageRecord struct {
	// RecordID is the internal ID of the water usage record
	RecordID int64 `db:"id"`

	// MunicipalityKey contains the AGS of the municipality the usage was recorded for
	MunicipalityKey string `db:"municipality"`

	// Date contains the timestamp indicating into which year the usage shall count
	Date pgtype.Timestamp `db:"date"`

	// Consumer contains the UUID of the consumer responsible for the usage. This value
	// may be null if no consumer was set
	Consumer pgtype.UUID `db:"consumer"`

	// UsageType contains a UUID pointing to the usage type for this record
	UsageType pgtype.UUID `db:"usage_type"`

	// RecordedAt contains a timestamp indicating the time at which the record
	// was written into the database
	RecordedAt pgtype.Timestamptz `db:"created_at"`

	// Amount contains the used water amount in cubic meters
	Amount float64 `db:"amount"`
}

// DbUsageType reflects a usage type stored in the WISdoM database. It may
// be converted to a json-encodable output by using the ToUsageType() function
// on the given instance
type DbUsageType struct {
	// ID contains the UUID used internally to reference the usage type
	ID pgtype.UUID `db:"id"`

	// Name contains the name or title of the usage type
	Title string `db:"name"`

	// Description contains the description of the usage type if one was set
	Description *string `db:"description"`

	// ExternalIdentifier contains the external identifier which is used in urls
	// and other frontend related actions
	ExternalIdentifier string `db:"external_identifier"`
}

// DbConsumer reflects a consumer stored in the database. It may
// be converted to a json-encodable output by using the ToConsumer() function
// on the given instance
type DbConsumer struct {
	// ID contains the database-generated UUID for this consumer
	ID pgtype.UUID `db:"id"`

	// Name contains the name of the consumer
	Name string `db:"name"`

	// Location contains the GeoJSON representation of the consumers location
	Location geojson.Geometry `db:"location"`

	// UsageType contains the UUID of the usage type that was assigned to the
	// consumer
	UsageType *pgtype.UUID `json:"usage_type"`

	// AdditionalProperties contains an optional key/value map allowing to add
	// additional properties to a consumer
	AdditionalProperties *pgtype.JSONB `json:"additional_properties"`
}
