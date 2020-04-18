package resources

import (
	"context"
	"net/http"
	"testing"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/felipecurvelo/weather-reporting-api/pkg/authorizer"
	"github.com/stretchr/testify/assert"
)

func TestWeatherFirstEndpoint_ReturnWelcomeMessage(t *testing.T) {
	ctx := context.Background()
	ctx = authorizer.NewContext(ctx, authorizer.MainAuth{})

	testServer := api.NewTestServer(ctx, t).RegisterResource(&Weather{})

	testServer.Test("GET", "/first_endpoint/").Now()
	statusCode, responseBody := testServer.GetResponse()

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Contains(t, responseBody, "This is the first endpoint working!")
}
