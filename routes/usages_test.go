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

func TestAllUsages_NoPageSize_NoPageNumber(t *testing.T) {
	var expectedHttpCode = http.StatusOK

	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/", AllUsages)

	request := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	router.ServeHTTP(recorder, request)

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

func TestAllUsages_NoPageSize_PageNumber(t *testing.T) {
	var expectedHttpCode = http.StatusOK

	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/", AllUsages)

	request := httptest.NewRequest("GET", "/?page=2", nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	router.ServeHTTP(recorder, request)

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

func TestAllUsages_PageSize_NoPageNumber(t *testing.T) {
	var expectedHttpCode = http.StatusOK
	var pageSize = 10

	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/", AllUsages)

	request := httptest.NewRequest("GET", fmt.Sprintf("/?page-size=%d", pageSize), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	router.ServeHTTP(recorder, request)

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

	if len(usageRecords) != pageSize {
		t.Errorf("Expected %d usage records, but got %d", DefaultPageSize, len(usageRecords))
	}
}

func TestAllUsages_PageSize_PageNumber(t *testing.T) {
	var expectedHttpCode = http.StatusOK
	var pageSize = 10
	var pageNumber = 2

	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/", AllUsages)

	request := httptest.NewRequest("GET", fmt.Sprintf("/?page-size=%d&page=%d", pageSize, pageNumber), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	router.ServeHTTP(recorder, request)

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

	if len(usageRecords) != pageSize {
		t.Errorf("Expected %d usage records, but got %d", DefaultPageSize, len(usageRecords))
	}
}

func TestAllUsages_PageSizeTooLarge_NoPageNumber(t *testing.T) {
	var expectedHttpCode = ErrPageTooLarge.Status
	var pageSize = MaxPageSize + 1

	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/", AllUsages)

	request := httptest.NewRequest("GET", fmt.Sprintf("/?page-size=%d", pageSize), nil)
	recorder := httptest.NewRecorder()

	_ = validator.NewValidator(apiContract).ForTest(t, recorder, request)
	router.ServeHTTP(recorder, request)

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
