package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/chi/v5"
	wisdomType "github.com/wisdom-oss/commonTypes/v2"
	wisdomMiddleware "github.com/wisdom-oss/microservice-middlewares/v4"

	"microservice/globals"
	"microservice/types"
)

const ARSKey = "ars"

var ErrEmptyARS = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: 400,
	Title:  "No ARS provided",
	Detail: "The request did not contain a ARS. Please check your request",
}

var ErrInvalidARS = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: 400,
	Title:  "Invalid ARS",
	Detail: "The ARS is not formatted correctly. Please check your request",
}

var ErrUnknownARS = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.5",
	Status: 404,
	Title:  "Unknown ARS",
	Detail: "The ARS is not present in the database",
}

func MunicipalUsages(w http.ResponseWriter, r *http.Request) {
	errorHandler := r.Context().Value(wisdomMiddleware.ErrorChannelName).(chan<- interface{})
	statusChannel := r.Context().Value(wisdomMiddleware.StatusChannelName).(<-chan bool)

	// check if the consumer id has been set correctly
	ars := strings.TrimSpace(chi.URLParam(r, ARSKey))
	if ars == "" {
		errorHandler <- ErrEmptyARS
		<-statusChannel
		return
	}

	if len(ars) != 12 {
		errorHandler <- ErrInvalidARS
		<-statusChannel
		return
	}

	// get the query returning all values
	query, err := globals.SqlQueries.Raw("check-municipal")
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	var municipalExists []bool
	err = pgxscan.Select(r.Context(), globals.Db, &municipalExists, query, ars)

	if len(municipalExists) == 0 || municipalExists[0] == false {
		errorHandler <- ErrUnknownARS
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

	filter, err := globals.SqlQueries.Raw("filter-municipal")
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

	rows, err := globals.Db.Query(r.Context(), query, ars, pageSize, offset)
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
