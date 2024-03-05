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
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	wisdomType "github.com/wisdom-oss/commonTypes/v2"
	wisdomMiddleware "github.com/wisdom-oss/microservice-middlewares/v4"

	"microservice/globals"
	"microservice/types"
)

const ConsumerIDKey = "consumer-id"

var ErrEmptyConsumerID = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: 400,
	Title:  "No Consumer ID provided",
	Detail: "The request did not contain a consumer id. Please check your request",
}

var ErrInvalidConsumerID = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: 400,
	Title:  "Invalid Consumer ID",
	Detail: "The consumer id is not formatted correctly. Please check your request",
}

var ErrUnknownConsumerID = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.5",
	Status: 404,
	Title:  "Unknown Consumer ID",
	Detail: "The consumer id is not present in the database",
}

func ConsumerUsages(w http.ResponseWriter, r *http.Request) {
	errorHandler := r.Context().Value(wisdomMiddleware.ErrorChannelName).(chan<- interface{})
	statusChannel := r.Context().Value(wisdomMiddleware.StatusChannelName).(<-chan bool)

	// check if the consumer id has been set correctly
	rawConsumerID := strings.TrimSpace(chi.URLParam(r, ConsumerIDKey))
	if rawConsumerID == "" {
		errorHandler <- ErrEmptyConsumerID
		<-statusChannel
		return
	}

	var consumerID pgtype.UUID
	err := consumerID.Scan(rawConsumerID)
	if err != nil {
		errorHandler <- ErrInvalidConsumerID
		<-statusChannel
		return
	}

	// get the query returning all values
	query, err := globals.SqlQueries.Raw("check-consumer")
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	var consumerExists []bool
	err = pgxscan.Select(r.Context(), globals.Db, &consumerExists, query, consumerID)

	if len(consumerExists) == 0 || consumerExists[0] == false {
		errorHandler <- ErrUnknownConsumerID
		<-statusChannel
		return
	}

	var outputFormat types.OutputFormat
	outputFormat.FromAcceptHeader(r.Header)

	var pageSize = DefaultPageSize
	var page = DefaultPage

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
	query, err = globals.SqlQueries.Raw("get-all")
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	filter, err := globals.SqlQueries.Raw("filter-consumer")
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	queryParts := strings.Split(query, "LIMIT")
	limit := queryParts[1]
	limit = strings.ReplaceAll(limit, "$2", "$3")
	limit = strings.ReplaceAll(limit, "$1", "$2")
	query = fmt.Sprintf("%s WHERE %s LIMIT %s", queryParts[0], filter, limit)

	rows, err := globals.Db.Query(r.Context(), query, consumerID, pageSize, offset)
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
