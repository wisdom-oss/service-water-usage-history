package db

import (
	"github.com/blockloop/scan/v2"
	"github.com/jackc/pgtype"
	"microservice/structs"
	"microservice/vars/globals"
	"microservice/vars/globals/connections"
)

// This file contains functions to allow the conversion from the database types
// to regular types which may be used for output and other things

// ToUsageRecord converts the database usage record into a usage record that may
// be used for api output. If an error occurs during the conversion the usage
// record will be nil and the error will be returned

func (r DbWaterUsageRecord) ToUsageRecord(usageTypes []DbUsageType) (*structs.UsageRecord, error) {
	// get the usage type that was associated with the usage record if not already pulled
	if usageTypes == nil {
		usageTypeRows, err := globals.Queries.Query(connections.DbConnection, "get-all-usage-types")
		if err != nil {
			return nil, err
		}
		// now try to parse the usage type
		err = scan.Rows(&usageTypes, usageTypeRows)
	}
	var dbUsageType DbUsageType
	for _, uT := range usageTypes {
		if uT.ID == r.UsageType {
			dbUsageType = uT
		}
	}

	// now check if the consumer is set on the record
	var consumerUUID *string
	if r.Consumer.Status != pgtype.Present {
		consumerUUID = nil
	} else {
		err := r.Consumer.AssignTo(&consumerUUID)
		// if the parsing failed pass it to the outside for handling
		if err != nil {
			return nil, err
		}
	}

	// now create the regular usage record and return it
	return &structs.UsageRecord{
		IncomingUsageRecord: structs.IncomingUsageRecord{
			Date:            r.Date.Time.Unix(),
			UsageType:       &dbUsageType.ExternalIdentifier,
			Consumer:        consumerUUID,
			Amount:          r.Amount,
			MunicipalityKey: r.MunicipalityKey,
		},
		RecordID:   r.RecordID,
		RecordedAt: r.RecordedAt.Time.Unix(),
	}, nil
}
