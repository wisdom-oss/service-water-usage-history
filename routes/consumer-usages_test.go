package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

// consumerTestRouter is the router used for the test cases as long as it is
// applicable to the test
var consumerTestRouter chi.Router

// consumerUsageTestMap contains the single test cases executed against the
// route
var consumerUsageTestMap = map[string]func(t *testing.T){
	"No Consumer ID":      noConsumerID,
	"Invalid Consumer ID": invalidConsumerID,
	"Valid Consumer ID":   validConsumerID,
	"Page Number":         consumerUsages_pageNumber,
	"Page Size":           consumerUsages_pageSize,
	"Page Size+Number":    consumerUsages_pageSizeAndNumber,
	"Page Size Too Large": consumerUsages_pageSizeTooLarge,
}

func TestConsumerUsages(t *testing.T) {

	consumerTestRouter = chi.NewRouter()
	consumerTestRouter.Use(middleware.ErrorHandler)
	consumerTestRouter.Get(fmt.Sprintf("/{%s}", ConsumerIDKey), ConsumerUsages)
	for testName, test := range consumerUsageTestMap {
		t.Run(testName, test)
	}

}

func noConsumerID(t *testing.T) {
	t.Parallel()
	expectedHttpCode := ErrEmptyConsumerID.Status

	// the path set here does not set the consumer-id since it is manually
	// added to the route context.
	// the path is only set to allow the correct detection of the path in the
	// api contract
	request := httptest.NewRequest("GET", "/abc", nil)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add(ConsumerIDKey, "")
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeContext))

	recorder := httptest.NewRecorder()
	_ = apiValidator.ForTest(t, recorder, request)

	consumerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ConsumerUsages(w, r)
	})

	mainHandler := middleware.ErrorHandler(consumerHandler)
	mainHandler.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		t.Errorf("Expected status code %d, but got %d", expectedHttpCode, recorder.Result().StatusCode)
	}

	// read the response and check that the default lengths have been used
	var apiError wisdomType.WISdoMError
	err := json.NewDecoder(recorder.Result().Body).Decode(&apiError)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if !apiError.Equals(ErrEmptyConsumerID) {
		t.Errorf("Expected error %v, but got %v", ErrEmptyConsumerID, apiError)
	}
}

func invalidConsumerID(t *testing.T) {
	t.Parallel()
	consumerID := "abc"
	expectedHttpCode := ErrInvalidConsumerID.Status

	request := httptest.NewRequest("GET", fmt.Sprintf("/%s", consumerID), nil)
	recorder := httptest.NewRecorder()

	_ = apiValidator.ForTest(t, recorder, request)
	consumerTestRouter.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		t.Errorf("Expected status code %d, but got %d", expectedHttpCode, recorder.Result().StatusCode)
	}

	// read the response and check that the default lengths have been used
	var apiError wisdomType.WISdoMError
	err := json.NewDecoder(recorder.Result().Body).Decode(&apiError)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if !apiError.Equals(ErrInvalidConsumerID) {
		t.Errorf("Expected error %v, but got %v", ErrInvalidConsumerID, apiError)
	}
}

func validConsumerID(t *testing.T) {
	t.Parallel()
	consumerID := "390dc645-c0a4-4cdf-8fbd-ab151f8c9687"
	expectedHttpCode := http.StatusOK

	request := httptest.NewRequest("GET", fmt.Sprintf("/%s", consumerID), nil)
	recorder := httptest.NewRecorder()

	_ = apiValidator.ForTest(t, recorder, request)
	consumerTestRouter.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		body, _ := io.ReadAll(recorder.Result().Body)
		var indentBuffer bytes.Buffer
		_ = json.Indent(&indentBuffer, body, "", "    ")
		t.Errorf("Expected status code %d, but got %d\n%s", expectedHttpCode, recorder.Result().StatusCode, indentBuffer.String())
	}

	// parse the response
	var usages []types.UsageRecord
	err := json.NewDecoder(recorder.Result().Body).Decode(&usages)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// check that the usages only contain the consumer id requested
	for _, usage := range usages {
		cID, _ := usage.Consumer.MarshalJSON()
		consumerIDString := strings.ReplaceAll(string(cID), `"`, "")
		if consumerIDString != consumerID {
			t.Errorf("Expected consumer ID %s, but got %s", consumerID, consumerIDString)
		}
	}

}

func consumerUsages_pageNumber(t *testing.T) {
	t.Parallel()
	consumerID := "390dc645-c0a4-4cdf-8fbd-ab151f8c9687"
	pageNumber := 1
	expectedHttpCode := http.StatusOK

	request := httptest.NewRequest("GET", fmt.Sprintf("/%s?page=%d", consumerID, pageNumber), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	consumerTestRouter.ServeHTTP(recorder, request)

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

	if len(usageRecords) > DefaultPageSize {
		t.Errorf("Expected max number of %d usage records, but got %d", DefaultPageSize, len(usageRecords))
	}
}

func consumerUsages_pageSize(t *testing.T) {
	t.Parallel()
	var expectedHttpCode = http.StatusOK

	consumerID := "390dc645-c0a4-4cdf-8fbd-ab151f8c9687"
	pageSize := 2

	request := httptest.NewRequest("GET", fmt.Sprintf("/%s?page-size=%d", consumerID, pageSize), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	consumerTestRouter.ServeHTTP(recorder, request)

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

	if len(usageRecords) > pageSize {
		t.Errorf("Expected max number of %d usage records, but got %d", DefaultPageSize, len(usageRecords))
	}
}

func consumerUsages_pageSizeAndNumber(t *testing.T) {
	t.Parallel()
	var expectedHttpCode = http.StatusOK

	consumerID := "390dc645-c0a4-4cdf-8fbd-ab151f8c9687"
	pageSize := 2
	pageNumber := 2

	request := httptest.NewRequest("GET", fmt.Sprintf("/%s?page-size=%d&page=%d", consumerID, pageSize, pageNumber), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	consumerTestRouter.ServeHTTP(recorder, request)

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

	if len(usageRecords) > pageSize {
		t.Errorf("Expected max number of %d usage records, but got %d", DefaultPageSize, len(usageRecords))
	}
}

func consumerUsages_pageSizeTooLarge(t *testing.T) {
	t.Parallel()
	var expectedHttpCode = ErrPageTooLarge.Status

	consumerID := "390dc645-c0a4-4cdf-8fbd-ab151f8c9687"
	pageSize := MaxPageSize + 1
	pageNumber := 1

	request := httptest.NewRequest("GET", fmt.Sprintf("/%s?page-size=%d&page=%d", consumerID, pageSize, pageNumber), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	consumerTestRouter.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		t.Errorf("Expected status code %d, but got %d", expectedHttpCode, recorder.Result().StatusCode)
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
