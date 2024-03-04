package routes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	middleware "github.com/wisdom-oss/microservice-middlewares/v4"
)

var consumerTestRouter chi.Router
var consumerUsageTestMap = map[string]func(t *testing.T){
	"No Consumer ID": noConsumerID,
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
	expectedHttpCode := http.StatusNotFound

	request := httptest.NewRequest("GET", "/", nil)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add(ConsumerIDKey, "")
	recorder := httptest.NewRecorder()

	_ = apiValidator.ForTest(t, recorder, request)
	consumerTestRouter.ServeHTTP(recorder, request)

	// Assert the response status code
	if recorder.Result().StatusCode != expectedHttpCode {
		t.Errorf("Expected status code %d, but got %d", expectedHttpCode, recorder.Result().StatusCode)
	}
}
