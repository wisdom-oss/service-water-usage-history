package routes

import (
	"encoding/json"
	"fmt"
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

var usageTestRouter chi.Router
var usageTestMap = map[string]func(t *testing.T){
	"noParameters":        noParameters,
	"pageNumber":          pageNumberOnly,
	"pageSize":            pageSizeOnly,
	"pageSize+pageNumber": pageSizeAndNumber,
	"pageSizeTooLarge":    pageSizeTooLarge,
	"output:JSON":         outputJSON,
	"output:CSV":          outputCSV,
	"output:CBOR":         outputCBOR,
}

func TestAllUsages(t *testing.T) {
	t.Parallel()
	usageTestRouter = chi.NewRouter()
	usageTestRouter.Use(middleware.ErrorHandler)
	usageTestRouter.Get("/", AllUsages)
	for testName, test := range usageTestMap {
		t.Run(testName, test)
	}
}

func noParameters(t *testing.T) {
	t.Parallel()
	var expectedHttpCode = http.StatusOK

	request := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	_ = apiValidator.ForTest(t, recorder, request)
	usageTestRouter.ServeHTTP(recorder, request)

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

	if len(usageRecords) != DefaultPageSize {
		t.Errorf("Expected %d usage records, but got %d", DefaultPageSize, len(usageRecords))
	}
}

func pageNumberOnly(t *testing.T) {
	t.Parallel()
	var expectedHttpCode = http.StatusOK
	request := httptest.NewRequest("GET", fmt.Sprintf("/?page=%d", pageNumber), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	usageTestRouter.ServeHTTP(recorder, request)

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

func pageSizeOnly(t *testing.T) {
	t.Parallel()
	var expectedHttpCode = http.StatusOK

	request := httptest.NewRequest("GET", fmt.Sprintf("/?page-size=%d", pageSize), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	usageTestRouter.ServeHTTP(recorder, request)

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
		t.Errorf("Expected max number of %d usage records, but got %d", pageSize, len(usageRecords))
	}
}

func pageSizeAndNumber(t *testing.T) {
	t.Parallel()
	var expectedHttpCode = http.StatusOK

	request := httptest.NewRequest("GET", fmt.Sprintf("/?page-size=%d&page=%d", pageSize, pageNumber), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	usageTestRouter.ServeHTTP(recorder, request)

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
		t.Errorf("Expected max number of %d usage records, but got %d", pageSize, len(usageRecords))
	}
}

func pageSizeTooLarge(t *testing.T) {
	t.Parallel()
	var expectedHttpCode = ErrPageTooLarge.Status
	var pageSize = MaxPageSize + 1

	request := httptest.NewRequest("GET", fmt.Sprintf("/?page-size=%d", pageSize), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	usageTestRouter.ServeHTTP(recorder, request)

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

func outputJSON(t *testing.T) {
	t.Parallel()
	var expectedHttpCode = http.StatusOK

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("Accept", "application/json")
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	usageTestRouter.ServeHTTP(recorder, request)

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

func outputCSV(t *testing.T) {
	t.Parallel()
	var expectedHttpCode = http.StatusOK

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("Accept", "text/csv")
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	usageTestRouter.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		t.Errorf("Expected status code %d, but got %d", expectedHttpCode, recorder.Result().StatusCode)
	}
}

func outputCBOR(t *testing.T) {
	t.Parallel()
	var expectedHttpCode = http.StatusOK

	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("Accept", "application/cbor")
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	usageTestRouter.ServeHTTP(recorder, request)

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
