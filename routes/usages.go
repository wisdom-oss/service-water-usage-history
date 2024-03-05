package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fxamacker/cbor/v2"
	"github.com/georgysavva/scany/v2/pgxscan"

	"microservice/globals"
	"microservice/types"

	wisdomMiddleware "github.com/wisdom-oss/microservice-middlewares/v4"
)

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
		err = encodeCSV(w, usages)
		break
	case types.CBOR:
		w.Header().Set("Content-Type", "application/cbor")
		err = cbor.NewEncoder(w).Encode(usages)
		break
	case types.JSON:
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(usages)
		break
	}

}
