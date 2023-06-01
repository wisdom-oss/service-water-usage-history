package routes

import (
	"github.com/blockloop/scan/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	requestErrors "microservice/request/error"
	"microservice/structs/db"
	"microservice/vars/globals"
	"microservice/vars/globals/connections"
	"net/http"
	"strconv"
)

var l = globals.HttpLogger

func GetAllUsages(w http.ResponseWriter, r *http.Request) {
	l.Warn().Msg("request for all usages received. response may take some time")
	// create a channel which the output data will be written to
	outputChan := make(chan render.Renderer, 50)
	// now start the database query in a goroutine, to allow the service to stay
	// responsive
	go func() {
		// get the usage types from the request context that have been inserted
		// by the middleware
		usageTypes := r.Context().Value("usageTypes").([]db.DbUsageType)
		// now get the query parameters indicating which results shall be returned
		pageStarts, pageStartSet := r.URL.Query()["start"]
		pageEnds, pageEndSet := r.URL.Query()["end"]
		// if the page start is not set, start at the first entry
		if !pageStartSet {
			pageStarts = append(pageStarts, "1")
		}
		// if the page end is not set, reject the request to protect the service
		// from failing
		if !pageEndSet {
			l.Warn().Str("query", "get-all-usages").Msg("no page end set. rejected")
			// build an error and put it into the output channel
			e, _ := requestErrors.GetRequestError("NO_PAGE_END")
			outputChan <- e
			// close the output channel to stop other messages to be written
			close(outputChan)
			return
		}
		// now get the first element of every one of the arrays
		pageStartStr := pageStarts[0]
		pageEndStr := pageEnds[0]
		// now try to convert the page bounds into integers
		pageStart, err := strconv.Atoi(pageStartStr)
		if err != nil {
			l.Warn().Err(err).Str("query", "get-all-usages").Msg("no page start set. assuming 1")
			pageStart = 1
		}
		pageEnd, err := strconv.Atoi(pageEndStr)
		if err != nil {
			l.Warn().Err(err).Str("query", "get-all-usages").Msg("page end is no string. rejecting")
			// build an error and put it into the output channel
			e, _ := requestErrors.GetRequestError("PAGE_BOUND_NOT_INT")
			outputChan <- e
			// close the output channel to stop other messages to be written
			close(outputChan)
			return
		}
		// now query the database for the usage records
		rows, err := globals.Queries.Query(connections.DbConnection, "get-all-usages",
			pageStart, pageEnd)
		if err != nil {
			l.Error().Err(err).Str("query", "get-all-usages").Msg("error during db query")
			// build an error and put it into the output channel
			e, _ := requestErrors.WrapInternalError(err)
			outputChan <- e
			// close the output channel to stop other messages to be written
			close(outputChan)
			return
		}
		// now iterate through the records and write them into the output channel
		for rows.Next() {
			// parse the row into a struct for easier handling
			var dbRecord db.DbWaterUsageRecord
			_ = rows.Scan(
				&dbRecord.RecordID,
				&dbRecord.MunicipalityKey,
				&dbRecord.Date,
				&dbRecord.Consumer,
				&dbRecord.UsageType,
				&dbRecord.RecordedAt,
				&dbRecord.Amount)
			// now convert the database-oriented struct into a returnable struct
			record, err := dbRecord.ToUsageRecord(usageTypes)
			if err != nil {
				l.Error().Err(err).Str("query", "get-all-usages").Msg("error during struct conversion")
				// build an error and put it into the output channel
				e, _ := requestErrors.WrapInternalError(err)
				outputChan <- e
				// close the output channel to stop other messages to be written
				close(outputChan)
				return
			}
			// send the record into the channel
			outputChan <- record
		}

		// since all rows have been handled close the channel
		close(outputChan)
	}()

	// since the data pulling is handled in a goroutine (async) start responding
	// with the objects from the output channel
	render.Respond(w, r, outputChan)
}

func GetConsumerUsages(w http.ResponseWriter, r *http.Request) {
	l.Warn().Msg("request for all usages received. response may take some time")
	// get the consumer id from the request
	consumerId := chi.URLParam(r, "consumerId")
	// create a channel which the output data will be written to
	outputChan := make(chan render.Renderer, 50)
	// now start the database query in a goroutine, to allow the service to stay
	// responsive
	go func() {
		l.Info().Msg("executing db query for all usage records")
		rows, err := globals.Queries.Query(connections.DbConnection, "get-consumers-usages", consumerId)
		// now check if the query returned an error
		if err != nil {
			l.Error().Err(err).Str("query", "get-consumers-usages").Msg("error during database query")
			// build an error and put it into the output channel
			e, _ := requestErrors.WrapInternalError(err)
			outputChan <- e
			// close the output channel to stop other messages to be written
			close(outputChan)
			return
		}

		// now start parsing the returned rows into objects
		var dbRecords []db.DbWaterUsageRecord
		err = scan.Rows(&dbRecords, rows)
		// now check if the query returned an error
		if err != nil {
			l.Error().Err(err).Msg("error during parsing of returned rows")
			// build an error and put it into the output channel
			e, _ := requestErrors.WrapInternalError(err)
			outputChan <- e
			// close the output channel to stop other messages to be written
			close(outputChan)
			return
		}
		// before converting the parsed rows access the usage types from the
		// request's context
		usageTypes := r.Context().Value("usageTypes").([]db.DbUsageType)
		// now iterate through the parsed rows
		for _, dbRecord := range dbRecords {
			// convert the usage record and send the result to the output
			record, err := dbRecord.ToUsageRecord(usageTypes)
			// now check if the conversion was successful
			if err != nil {
				l.Error().Err(err).Msg("error during conversion to output")
				// build an error and put it into the output channel
				e, _ := requestErrors.WrapInternalError(err)
				outputChan <- e
				break
			}
			// since the conversion was successful, pass the converted record
			// into the output channel
			outputChan <- record
		}
		// close the output channel to stop other messages to be written
		close(outputChan)
	}()
	// start responding with the objects in the output channel
	render.Respond(w, r, outputChan)
}
