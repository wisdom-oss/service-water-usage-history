package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"microservice/structs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/stretchr/testify/assert"
)

func _paged_usages(t *testing.T) {
	t.Parallel()
	t.Run("Defaults", _pu_defaults)
	t.Run("Pagination", _pu_pages)
	t.Run("Sizing", _pu_page_size)
}

func _pu_defaults(t *testing.T) {
	t.Parallel()
	expectedEntries := 10000

	req := httptest.NewRequest("GET", routePrefix+"/", nil)
	res := httptest.NewRecorder()

	r.Handler().ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	err := openapi3filter.ValidateResponse(context.Background(), generateValidationData(t, req, res))
	if err != nil {
		t.Fail()
		t.Log(err)
	}

	var entries []structs.UsageRecord
	err = json.NewDecoder(res.Result().Body).Decode(&entries)
	assert.NoError(t, err)
	assert.Len(t, entries, expectedEntries)
}

func _pu_pages(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", routePrefix+"/", nil)
	res := httptest.NewRecorder()

	r.Handler().ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	err := openapi3filter.ValidateResponse(context.Background(), generateValidationData(t, req, res))
	if err != nil {
		t.Fail()
		t.Log(err)
	}

	var page1 []structs.UsageRecord
	err = json.NewDecoder(res.Result().Body).Decode(&page1)
	assert.NoError(t, err)

	req = httptest.NewRequest("GET", routePrefix+"/?page=2", nil)
	res = httptest.NewRecorder()
	r.Handler().ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	err = openapi3filter.ValidateResponse(context.Background(), generateValidationData(t, req, res))
	if err != nil {
		t.Fail()
		t.Log(err)
	}

	var page2 []structs.UsageRecord
	err = json.NewDecoder(res.Result().Body).Decode(&page2)
	assert.NoError(t, err)

	assert.NotElementsMatch(t, page1, page2)
}

func _pu_page_size(t *testing.T) {
	t.Parallel()
	expectedEntries := 1000

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/?pageSize=%d", routePrefix, expectedEntries), nil)
	res := httptest.NewRecorder()

	r.Handler().ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	err := openapi3filter.ValidateResponse(context.Background(), generateValidationData(t, req, res))
	if err != nil {
		t.Fail()
		t.Log(err)
	}

	var entries []structs.UsageRecord
	err = json.NewDecoder(res.Result().Body).Decode(&entries)
	assert.NoError(t, err)
	assert.Len(t, entries, expectedEntries)
}
