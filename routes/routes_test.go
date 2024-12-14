package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

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
	r.GET("/", PagedUsages)

	t.Run("Paged_Usages", _paged_usages)
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
