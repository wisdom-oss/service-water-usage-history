package routes

import (
	"context"
	"encoding/json"
	"fmt"
	apiErrors "microservice/internal/errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/stretchr/testify/assert"
	"github.com/wisdom-oss/common-go/v3/types"
)

func _municipal_usages(t *testing.T) {

	t.Run("Empty_ARS", _mu_empty_ars)
	t.Run("Invalid_ARS", _mu_invalid_ars)
	t.Run("Valid_Request", _mu_valid_request)
}

func _mu_empty_ars(t *testing.T) {
	apiPath := "municipal"
	pathParameter := ""
	expectedError := apiErrors.ErrEmptyARS

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/%s/%s", routePrefix, apiPath, pathParameter), nil)
	res := httptest.NewRecorder()

	r.Handler().ServeHTTP(res, req)
	assert.Equal(t, int(expectedError.Status), res.Code)

	err := openapi3filter.ValidateResponse(context.Background(), generateValidationData(t, req, res))
	if err != nil {
		t.Fail()
		t.Log(err)
	}

	var receivedError types.ServiceError
	err = json.NewDecoder(res.Body).Decode(&receivedError)
	assert.NoError(t, err)
	if t.Failed() {
		t.FailNow()
	}

	assert.True(t, receivedError.Equals(expectedError))
}

func _mu_invalid_ars(t *testing.T) {
	apiPath := "municipal"
	pathParameter := `031515401020`[0:10]
	expectedError := apiErrors.ErrInvalidARS

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/%s/%s", routePrefix, apiPath, pathParameter), nil)
	res := httptest.NewRecorder()

	r.Handler().ServeHTTP(res, req)
	assert.Equal(t, int(expectedError.Status), res.Code)

	err := openapi3filter.ValidateResponse(context.Background(), generateValidationData(t, req, res))
	if err != nil {
		t.Fail()
		t.Log(err)
	}

	var receivedError types.ServiceError
	err = json.NewDecoder(res.Body).Decode(&receivedError)
	assert.NoError(t, err)
	if t.Failed() {
		t.FailNow()
	}

	assert.True(t, receivedError.Equals(expectedError))
}

func _mu_valid_request(t *testing.T) {
	apiPath := "municipal"
	pathParameter := `031515401020`

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/%s/%s", routePrefix, apiPath, pathParameter), nil)
	res := httptest.NewRecorder()

	r.Handler().ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	err := openapi3filter.ValidateResponse(context.Background(), generateValidationData(t, req, res))
	if err != nil {
		t.Fail()
		t.Log(err)
	}
}
