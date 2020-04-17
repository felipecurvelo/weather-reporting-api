package resources

import (
	"net/http"
	"testing"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestWeatherFirstEndpoint_ReturnWelcomeMessage(t *testing.T) {
	testServer := api.NewTestServer().RegisterResource(&Weather{})

	responseBody, statusCode, err := testServer.CallEndpoint("GET", "/first_endpoint/")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "Welcome! This is the first endpoint working!", responseBody)
}
