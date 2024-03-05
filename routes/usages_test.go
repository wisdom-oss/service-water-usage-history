package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	wisdomType "github.com/wisdom-oss/commonTypes/v2"
	middleware "github.com/wisdom-oss/microservice-middlewares/v4"
	validator "openapi.tanna.dev/go/validator/openapi3"

	"microservice/types"
)

var usageTestRouter chi.Router
var usageTestMap = map[string]func(t *testing.T){
	"No Parameters":        noParameters,
	"Page Number":          pageNumberOnly,
	"Page Size":            pageSizeOnly,
	"Page Size and Number": pageSizeAndNumber,
	"Page Size Too Large":  pageSizeTooLarge,
}

func TestAllUsages(t *testing.T) {
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
	request := httptest.NewRequest("GET", "/?page=2", nil)
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
	var pageSize = 10

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
		t.Errorf("Expected max number of %d usage records, but got %d", DefaultPageSize, len(usageRecords))
	}
}

func pageSizeAndNumber(t *testing.T) {
	t.Parallel()
	var expectedHttpCode = http.StatusOK
	var pageSize = 10
	var pageNumber = 2

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
		t.Errorf("Expected max number of %d usage records, but got %d", DefaultPageSize, len(usageRecords))
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
