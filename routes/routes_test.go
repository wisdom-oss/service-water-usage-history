package routes

import (
	"context"
	"encoding/json"
	apiErrors "microservice/internal/errors"
	routeUtils "microservice/routes/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/wisdom-oss/common-go/v3/types"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

var r *gin.Engine
var openapi *openapi3.T

// routePrefix is used to allow the correct resolving of a route for the
// openapi validation library
const routePrefix = `http://localhost:8000`

func TestRoutes(t *testing.T) {

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("../openapi.yaml")
	if err != nil {
		panic(err)
	}

	err = doc.Validate(context.Background())
	if err != nil {
		t.Log("Invalid OpenAPI document found")
		t.FailNow()
	}

	openapi = doc

	r = gin.New()
	r.Use(routeUtils.ReadPageSettings)
	r.GET("/", PagedUsages)
	r.GET("/consumer/*consumerID", ConsumerUsages)
	r.GET("/municipal/*ars", MunicipalUsages)
	r.GET("/type/*usageTypeID", TypedUsages)

	t.Run("Paged_Usages", _paged_usages)
	t.Run("Consumer_Usages", _consumer_usages)
	t.Run("Municipal_Usages", _municipal_usages)
	t.Run("Typed_Usages", _typed_usages)
	t.Run("Page_Settings", _page_settings)
}

func generateValidationData(t *testing.T, req *http.Request, res *httptest.ResponseRecorder) *openapi3filter.ResponseValidationInput {
	router, err := gorillamux.NewRouter(openapi)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	route, pathParams, err := router.FindRoute(req)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	requestValidation := &openapi3filter.RequestValidationInput{
		Request:    req,
		PathParams: pathParams,
		Route:      route,
	}

	responseValidation := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidation,
		Status:                 res.Code,
		Header:                 res.Header(),
	}

	responseValidation.SetBodyBytes(res.Body.Bytes())
	return responseValidation
}

func _page_settings(t *testing.T) {
	t.Run("Page_To_Small", func(t *testing.T) {
		expectedError := apiErrors.ErrInvalidPageSettings

		t.Parallel()
		req := httptest.NewRequest("GET", routePrefix+"/?page=-1", nil)
		res := httptest.NewRecorder()

		r.Handler().ServeHTTP(res, req)
		assert.Equal(t, int(expectedError.Status), res.Code)

		var receivedError types.ServiceError
		err := json.NewDecoder(res.Body).Decode(&receivedError)
		assert.NoError(t, err)
		if t.Failed() {
			t.FailNow()
		}

		assert.True(t, receivedError.Equals(expectedError))

	})

	t.Run("Page_Size_Negative", func(t *testing.T) {
		expectedError := apiErrors.ErrInvalidPageSettings

		t.Parallel()
		req := httptest.NewRequest("GET", routePrefix+"/?pageSize=-1", nil)
		res := httptest.NewRecorder()

		r.Handler().ServeHTTP(res, req)
		assert.Equal(t, int(expectedError.Status), res.Code)

		var receivedError types.ServiceError
		err := json.NewDecoder(res.Body).Decode(&receivedError)
		assert.NoError(t, err)
		if t.Failed() {
			t.FailNow()
		}

		assert.True(t, receivedError.Equals(expectedError))

	})

	t.Run("Page_Size_Too_High", func(t *testing.T) {
		expectedError := apiErrors.ErrInvalidPageSettings

		t.Parallel()
		req := httptest.NewRequest("GET", routePrefix+"/?pageSize=100001", nil)
		res := httptest.NewRecorder()

		r.Handler().ServeHTTP(res, req)
		assert.Equal(t, int(expectedError.Status), res.Code)

		var receivedError types.ServiceError
		err := json.NewDecoder(res.Body).Decode(&receivedError)
		assert.NoError(t, err)
		if t.Failed() {
			t.FailNow()
		}

		assert.True(t, receivedError.Equals(expectedError))

	})
}
