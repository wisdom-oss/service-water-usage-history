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

func _consumer_usages(t *testing.T) {
	t.Parallel()
	t.Run("Empty_Consumer_ID", _cu_empty_consumer_id)
	t.Run("Invalid_Consumer_ID", _cu_invalid_consumer_id)
	t.Run("Valid_Request", _cu_valid_request)
}

func _cu_empty_consumer_id(t *testing.T) {
	apiPath := "consumer"
	pathParameter := ""
	expectedError := apiErrors.ErrEmptyConsumerID

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

func _cu_invalid_consumer_id(t *testing.T) {
	apiPath := "consumer"
	pathParameter := randstr.Base64(rand.Intn(64))
	expectedError := apiErrors.ErrInvalidConsumerID

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

func _cu_valid_request(t *testing.T) {
	apiPath := "consumer"
	pathParameter := `390dc645-c0a4-4cdf-8fbd-ab151f8c9687`

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
