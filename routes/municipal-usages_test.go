package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/go-chi/chi/v5"
	wisdomType "github.com/wisdom-oss/commonTypes/v2"
	middleware "github.com/wisdom-oss/microservice-middlewares/v4"
	validator "openapi.tanna.dev/go/validator/openapi3"

	"microservice/types"
)

// This file contains all test cases executed against the ConsumerUsages
// handler.
// Even though the handler is directly derived from the AllUsages route, the
// tests for the pagination handling will be repeated to ensure a working code
// on all routes

// municipalTestRouter is the router used for the test cases as long as it is
// applicable to the test
var municipalTestRouter chi.Router

const ars = "032565409007"

// municipalUsageTestMap contains the single test cases executed against the
// route
var municipalUsageTestMap = map[string]func(t *testing.T){
	"noARS":               noARS,
	"invalidARS":          invalidARS,
	"validARS":            validARS,
	"pageNumber":          municipalUsages_pageNumber,
	"pageSize":            municipalUsages_pageSize,
	"pageSize+pageNumber": municipalUsages_pageSizeAndNumber,
	"pageSizeTooLarge":    municipalUsages_pageSizeTooLarge,
	"output:JSON":         municipalUsages_outputJSON,
	"output:CSV":          municipalUsages_outputCSV,
	"output:CBOR":         municipalUsages_outputCBOR,
}

func TestMunicipalUsages(t *testing.T) {

	municipalTestRouter = chi.NewRouter()
	municipalTestRouter.Use(middleware.ErrorHandler)
	municipalTestRouter.Get(fmt.Sprintf("/municipal/{%s}", ARSKey), MunicipalUsages)
	for testName, test := range municipalUsageTestMap {
		t.Run(testName, test)
	}

}

func noARS(t *testing.T) {

	expectedHttpCode := ErrEmptyConsumerID.Status

	// the path set here does not set the municipal-id since it is manually
	// added to the route context.
	// the path is only set to allow the correct detection of the path in the
	// api contract
	request := httptest.NewRequest("GET", "/municipal/abc", nil)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add(ARSKey, "")
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeContext))

	recorder := httptest.NewRecorder()
	_ = apiValidator.ForTest(t, recorder, request)

	municipalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MunicipalUsages(w, r)
	})

	mainHandler := middleware.ErrorHandler(municipalHandler)
	mainHandler.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		body, _ := io.ReadAll(recorder.Result().Body)
		var indentBuffer bytes.Buffer
		_ = json.Indent(&indentBuffer, body, "", "    ")
		t.Fatalf("Expected status code %d, but got %d\n%s", expectedHttpCode, recorder.Result().StatusCode, indentBuffer.String())
	}

	// read the response and check that the default lengths have been used
	var apiError wisdomType.WISdoMError
	err := json.NewDecoder(recorder.Result().Body).Decode(&apiError)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if !apiError.Equals(ErrEmptyARS) {
		t.Errorf("Expected error %v, but got %v", ErrEmptyConsumerID, apiError)
	}
}

func invalidARS(t *testing.T) {

	ars := "abc"
	expectedHttpCode := ErrInvalidConsumerID.Status

	request := httptest.NewRequest("GET", fmt.Sprintf("/municipal/%s", ars), nil)
	recorder := httptest.NewRecorder()

	_ = apiValidator.ForTest(t, recorder, request)
	municipalTestRouter.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		body, _ := io.ReadAll(recorder.Result().Body)
		var indentBuffer bytes.Buffer
		_ = json.Indent(&indentBuffer, body, "", "    ")
		t.Fatalf("Expected status code %d, but got %d\n%s", expectedHttpCode, recorder.Result().StatusCode, indentBuffer.String())
	}

	// read the response and check that the default lengths have been used
	var apiError wisdomType.WISdoMError
	err := json.NewDecoder(recorder.Result().Body).Decode(&apiError)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if !apiError.Equals(ErrInvalidARS) {
		t.Errorf("Expected error %v, but got %v", ErrInvalidARS, apiError)
	}
}

func validARS(t *testing.T) {

	expectedHttpCode := http.StatusOK

	request := httptest.NewRequest("GET", fmt.Sprintf("/municipal/%s", ars), nil)
	recorder := httptest.NewRecorder()

	_ = apiValidator.ForTest(t, recorder, request)
	municipalTestRouter.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		body, _ := io.ReadAll(recorder.Result().Body)
		var indentBuffer bytes.Buffer
		_ = json.Indent(&indentBuffer, body, "", "    ")
		t.Fatalf("Expected status code %d, but got %d\n%s", expectedHttpCode, recorder.Result().StatusCode, indentBuffer.String())
	}

	// parse the response
	var usages []types.UsageRecord
	err := json.NewDecoder(recorder.Result().Body).Decode(&usages)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// check that the usages only contain the municipal id requested
	for _, usage := range usages {
		municipalID := usage.Municipality.String
		if municipalID != ars {
			t.Errorf("Expected municipal ID %s, but got %s", ars, municipalID)
		}
	}

}

func municipalUsages_pageNumber(t *testing.T) {

	expectedHttpCode := http.StatusOK

	request := httptest.NewRequest("GET", fmt.Sprintf("/municipal/%s?page=%d", ars, pageNumber), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	municipalTestRouter.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		body, _ := io.ReadAll(recorder.Result().Body)
		var indentBuffer bytes.Buffer
		_ = json.Indent(&indentBuffer, body, "", "    ")
		t.Fatalf("Expected status code %d, but got %d\n%s", expectedHttpCode, recorder.Result().StatusCode, indentBuffer.String())
	}

	// read the response and check that the default lengths have been used
	var usageRecords []types.UsageRecord
	err := json.NewDecoder(recorder.Result().Body).Decode(&usageRecords)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(usageRecords) > DefaultPageSize {
		t.Errorf("Expected max number of %d usage records, but got %d", DefaultPageSize, len(usageRecords))
	}
}

func municipalUsages_pageSize(t *testing.T) {

	var expectedHttpCode = http.StatusOK

	request := httptest.NewRequest("GET", fmt.Sprintf("/municipal/%s?page-size=%d", ars, pageSize), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	municipalTestRouter.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		body, _ := io.ReadAll(recorder.Result().Body)
		var indentBuffer bytes.Buffer
		_ = json.Indent(&indentBuffer, body, "", "    ")
		t.Fatalf("Expected status code %d, but got %d\n%s", expectedHttpCode, recorder.Result().StatusCode, indentBuffer.String())
	}
}

func municipalUsages_pageSizeAndNumber(t *testing.T) {

	var expectedHttpCode = http.StatusOK

	request := httptest.NewRequest("GET", fmt.Sprintf("/municipal/%s?page-size=%d&page=%d", ars, pageSize, pageNumber), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	municipalTestRouter.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		body, _ := io.ReadAll(recorder.Result().Body)
		var indentBuffer bytes.Buffer
		_ = json.Indent(&indentBuffer, body, "", "    ")
		t.Fatalf("Expected status code %d, but got %d\n%s", expectedHttpCode, recorder.Result().StatusCode, indentBuffer.String())
	}

	// read the response and check that the default lengths have been used
	var usageRecords []types.UsageRecord
	err := json.NewDecoder(recorder.Result().Body).Decode(&usageRecords)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(usageRecords) > pageSize {
		t.Errorf("Expected max number of %d usage records, but got %d", DefaultPageSize, len(usageRecords))
	}
}

func municipalUsages_pageSizeTooLarge(t *testing.T) {

	var expectedHttpCode = ErrPageTooLarge.Status

	pageSize := MaxPageSize + 1

	request := httptest.NewRequest("GET", fmt.Sprintf("/municipal/%s?page-size=%d&page=%d", ars, pageSize, pageNumber), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	municipalTestRouter.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		body, _ := io.ReadAll(recorder.Result().Body)
		var indentBuffer bytes.Buffer
		_ = json.Indent(&indentBuffer, body, "", "    ")
		t.Fatalf("Expected status code %d, but got %d\n%s", expectedHttpCode, recorder.Result().StatusCode, indentBuffer.String())
	}

	// read the response and check that the default lengths have been used
	var apiError wisdomType.WISdoMError
	err := json.NewDecoder(recorder.Result().Body).Decode(&apiError)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if !apiError.Equals(ErrPageTooLarge) {
		t.Errorf("Expected error %v, but got %v", ErrPageTooLarge, apiError)
	}
}

func municipalUsages_outputJSON(t *testing.T) {

	var expectedHttpCode = http.StatusOK

	request := httptest.NewRequest("GET", fmt.Sprintf("/municipal/%s", ars), nil)
	request.Header.Set("Accept", "application/json")
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	municipalTestRouter.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		t.Errorf("Expected status code %d, but got %d", expectedHttpCode, recorder.Result().StatusCode)
	}

	// read the response and check that the default lengths have been used
	var usageRecords []types.UsageRecord
	err := json.NewDecoder(recorder.Result().Body).Decode(&usageRecords)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
}

func municipalUsages_outputCSV(t *testing.T) {

	var expectedHttpCode = http.StatusOK

	request := httptest.NewRequest("GET", fmt.Sprintf("/municipal/%s", ars), nil)
	request.Header.Set("Accept", "text/csv")
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	municipalTestRouter.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		t.Errorf("Expected status code %d, but got %d", expectedHttpCode, recorder.Result().StatusCode)
	}
}

func municipalUsages_outputCBOR(t *testing.T) {

	var expectedHttpCode = http.StatusOK

	request := httptest.NewRequest("GET", fmt.Sprintf("/municipal/%s", ars), nil)
	request.Header.Set("Accept", "application/cbor")
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	municipalTestRouter.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		t.Errorf("Expected status code %d, but got %d", expectedHttpCode, recorder.Result().StatusCode)
	}

	// read the response and check that the default lengths have been used
	var usageRecords []types.UsageRecord
	err := cbor.NewDecoder(recorder.Result().Body).Decode(&usageRecords)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
}
