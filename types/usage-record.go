package types

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UsageRecord struct {
	Timestamp    pgtype.Timestamptz `json:"timestamp" db:"time" cbor:"1"`
	Amount       pgtype.Float8      `json:"amount" db:"amount" cbor:"2"`
	UsageType    *pgtype.UUID       `json:"usageType,omitempty" db:"usage_type" cbor:"3,omitempty"`
	Consumer     *pgtype.UUID       `json:"consumer,omitempty" db:"consumer" cbor:"4,omitempty"`
	Municipality pgtype.Text        `json:"municipality,omitempty" db:"municipality" cbor:"5,omitempty"`
}
