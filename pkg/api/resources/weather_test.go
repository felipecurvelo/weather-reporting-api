package resources

import (
	"context"
	"net/http"
	"testing"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/felipecurvelo/weather-reporting-api/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestWeatherFirstEndpoint_ReturnWelcomeMessage(t *testing.T) {
	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{
		Name: "felipe",
	})

	testServer := api.NewTestServer(ctx, t).RegisterResource(&Weather{})

	testServer.Test("GET", "/first_endpoint/").Now()
	statusCode, responseBody := testServer.GetResponse()

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "\"Welcome felipe! This is the first endpoint working!\"", responseBody)
}
