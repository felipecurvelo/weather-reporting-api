package resources

import (
	"net/http"
	"testing"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestWeatherFirstEndpoint_ReturnWelcomeMessage(t *testing.T) {
	testServer := api.NewTestServer(t).RegisterResource(&Weather{})

	testServer.Test("GET", "/first_endpoint/").Now()
	statusCode, responseBody := testServer.GetResponse()

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "\"Welcome! This is the first endpoint working!\"", responseBody)
}
