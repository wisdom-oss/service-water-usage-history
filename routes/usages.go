package routes

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/georgysavva/scany/v2/pgxscan"
	wisdomType "github.com/wisdom-oss/commonTypes/v2"

	"microservice/globals"
	"microservice/types"

	wisdomMiddleware "github.com/wisdom-oss/microservice-middlewares/v4"
)

var ErrPageTooLarge = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.14",
	Status: 413,
	Title:  "Page Size Too Large",
	Detail: fmt.Sprintf("Due to limitations on the system side, the selected page size is too large too handle. Please select a value smaller than %d", MaxPageSize),
}

const DefaultPageSize = 10000
const DefaultPage = 1
const MaxPageSize = 2500000

func AllUsages(w http.ResponseWriter, r *http.Request) {
	errorHandler := r.Context().Value(wisdomMiddleware.ErrorChannelName).(chan<- interface{})
	statusChannel := r.Context().Value(wisdomMiddleware.StatusChannelName).(<-chan bool)
	var outputFormat types.OutputFormat
	outputFormat.FromAcceptHeader(r.Header)

	var pageSize = DefaultPageSize
	var page = DefaultPage

	var err error

	pageSizeStrings, pageSizeSet := r.URL.Query()["page-size"]
	if pageSizeSet {
		pageSize, err = strconv.Atoi(pageSizeStrings[0])
		if err != nil {
			pageSize = DefaultPageSize
		}
	}

	if pageSize > MaxPageSize {
		errorHandler <- ErrPageTooLarge
		<-statusChannel
		return
	}

	pageStrings, pageSet := r.URL.Query()["page"]
	if pageSet {
		page, err = strconv.Atoi(pageStrings[0])
		if err != nil {
			page = DefaultPage
		}
	}

	offset := pageSize * (page - 1)

	// get the query returning all values
	query, err := globals.SqlQueries.Raw("get-all")
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	rows, err := globals.Db.Query(r.Context(), query, pageSize, offset)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	var usages []types.UsageRecord
	err = pgxscan.ScanAll(&usages, rows)

	if len(usages) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	switch outputFormat {
	case types.CSV:
		w.Header().Set("Content-Type", "text/csv")
		csvWriter := csv.NewWriter(w)
		_ = csvWriter.Write([]string{"timestamp", "amount", "usage-type", "consumer", "municipality"})
		for _, usageRecord := range usages {
			csvWriter.Flush()
			timestampString := usageRecord.Timestamp.Time.Format(time.RFC3339)
			value := fmt.Sprintf("%f", usageRecord.Amount.Float64)
			consumerID, _ := usageRecord.Consumer.MarshalJSON()
			consumerIDString := strings.ReplaceAll(string(consumerID), `"`, "")
			usageTypeID, _ := usageRecord.UsageType.MarshalJSON()
			usageTypeIDString := strings.ReplaceAll(string(usageTypeID), `"`, "")
			municipality := usageRecord.Municipality.String
			err = csvWriter.Write([]string{timestampString, value, usageTypeIDString, consumerIDString, municipality})
			if err != nil {
				errorHandler <- err
				<-statusChannel
				return
			}
			csvWriter.Flush()
		}

	case types.CBOR:
		w.Header().Set("Content-Type", "application/cbor")
		err = cbor.NewEncoder(w).Encode(usages)
		break
	case types.JSON:
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(usages)
	}

}
