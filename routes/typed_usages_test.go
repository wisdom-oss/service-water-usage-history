package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	apiErrors "microservice/internal/errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
	"github.com/wisdom-oss/common-go/v3/types"
)

func _typed_usages(t *testing.T) {
	t.Parallel()
	t.Run("Empty_Usage_Type_ID", _tu_empty_usage_type_id)
	t.Run("Invalid_Usage_Type_ID", _tu_invalid_usage_type_id)
	t.Run("Valid_Request", _tu_valid_request)
}

func _tu_empty_usage_type_id(t *testing.T) {
	apiPath := "type"
	pathParameter := ""
	expectedError := apiErrors.ErrEmptyUsageTypeID

	t.Parallel()

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

func _tu_invalid_usage_type_id(t *testing.T) {
	apiPath := "type"
	pathParameter := randstr.Base64(rand.Intn(64))
	expectedError := apiErrors.ErrInvalidUsageTypeID

	t.Parallel()

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

func _tu_valid_request(t *testing.T) {
	apiPath := "type"
	pathParameter := `d9e1dd0b-c25e-45c2-be40-c274b4845ad1`

	t.Parallel()

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
