package resources

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/felipecurvelo/weather-reporting-api/pkg/authorizer"
	"github.com/stretchr/testify/assert"
)

func TestWeatherFirstEndpoint_ReturnWelcomeMessage(t *testing.T) {
	ctx := authorizer.NewContext(context.Background(), authorizer.NewAuthMock())

	testServer := api.NewTestServer(ctx, t).
		RegisterResource(&Auth{}).
		RegisterResource(&Weather{})

	authRequestBody := `
		{
			"name": "kirang",
			"password": "secret"
		}
	`

	// Call the auth endpoint to get the auth token
	testServer.Test("POST", "/auth/").
		WithBody(authRequestBody).
		Now()

	statusCode, responseBody := testServer.GetResponse()
	assert.Equal(t, http.StatusOK, statusCode)

	// Parse the auth response
	var tokenObj map[string]interface{}
	err := json.Unmarshal([]byte(responseBody), &tokenObj)
	assert.NoError(t, err)

	// Uses the parsed auth as authorization header
	testServer.Test("GET", "/first_endpoint/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Contains(t, responseBody, "This is the first endpoint working!")
}
