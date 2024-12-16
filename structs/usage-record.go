package structs

import "github.com/jackc/pgx/v5/pgtype"

type UsageRecord struct {
	Time       pgtype.Timestamptz `json:"time" db:"time"`
	Amount     float64            `json:"amount" db:"amount"`
	UsageType  *string            `json:"usageType" db:"usage_type"`
	ConsumerID *string            `json:"consumerID" db:"consumer"`
	ARS        *string            `json:"ars" db:"municipality"`
}
