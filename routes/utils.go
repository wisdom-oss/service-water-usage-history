package routes

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
	"time"

	"microservice/types"
)

func encodeCSV(w http.ResponseWriter, usages []types.UsageRecord) error {
	enc := csv.NewWriter(w)
	defer enc.Flush()

	headers := []string{"timestamp", "amount", "usage-type", "consumer", "municipality"}
	if err := enc.Write(headers); err != nil {
		return err
	}

	for _, usage := range usages {
		timestamp := usage.Timestamp.Time.Format(time.RFC3339)
		amount := fmt.Sprintf("%f", usage.Amount.Float64)
		ars := usage.Municipality.String

		var usageType string
		if usage.UsageType != nil {
			usageTypeBytes, _ := usage.UsageType.MarshalJSON()
			usageType = string(usageTypeBytes)
			usageType = strings.ReplaceAll(usageType, `"`, ``)
		}

		var consumer string
		if usage.Consumer != nil {
			consumerBytes, _ := usage.Consumer.MarshalJSON()
			consumer = string(consumerBytes)
			consumer = strings.ReplaceAll(consumer, `"`, ``)
		}

		err := enc.Write([]string{timestamp, amount, usageType, consumer, ars})
		if err != nil {
			return err
		}
		enc.Flush()
	}

	return nil
}
