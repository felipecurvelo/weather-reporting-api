package resources

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/felipecurvelo/weather-reporting-api/pkg/authorizer"
	"github.com/felipecurvelo/weather-reporting-api/pkg/weathermanager"
	"github.com/stretchr/testify/assert"
)

func TestWeatherSave_ReturnOK(t *testing.T) {
	ctx := authorizer.NewContext(context.Background(), authorizer.NewAuthMock())
	ctx = weathermanager.NewContext(ctx, weathermanager.New())

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

	saveRequestBody := `
		{
			"city": "vancouver",
			"weather": [{
				"date": "2020-04-18",
				"temperature": 15
			}]
		}
	`

	// Uses the parsed auth as authorization header to save weather
	testServer.Test("POST", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(saveRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "{\"message\":\"The weather was saved succesfully!\"}", responseBody)
}

func TestWeatherSave_AndRetrieve_ReturnVancouver(t *testing.T) {
	ctx := authorizer.NewContext(context.Background(), authorizer.NewAuthMock())
	ctx = weathermanager.NewContext(ctx, weathermanager.New())

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

	saveRequestBody := `
		{
			"city": "vancouver",
			"weather": [{
				"date": "2020-04-18",
				"temperature": 15
			}]
		}
	`

	// Uses the parsed auth as authorization header to save weather
	testServer.Test("POST", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(saveRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "{\"message\":\"The weather was saved succesfully!\"}", responseBody)

	getRequestBody := `
		{
			"city": "vancouver",
			"initial_date": "2020-04-01",
			"end_date": "2020-04-30"
		}
	`

	// Gets the saved weather
	testServer.Test("GET", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(getRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Contains(t, responseBody, "vancouver")
}

func TestWeatherSave_AndDelete_ReturnEmpty(t *testing.T) {
	ctx := authorizer.NewContext(context.Background(), authorizer.NewAuthMock())
	ctx = weathermanager.NewContext(ctx, weathermanager.New())

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

	saveRequestBody := `
		{
			"city": "vancouver",
			"weather": [{
				"date": "2020-04-18",
				"temperature": 15
			},
			{
				"date": "2020-05-18",
				"temperature": 16
			}]
		}
	`

	// Uses the parsed auth as authorization header to save weather
	testServer.Test("POST", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(saveRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "{\"message\":\"The weather was saved succesfully!\"}", responseBody)

	getRequestBody := `
		{
			"city": "vancouver",
			"initial_date": "2020-04-01",
			"end_date": "2020-04-30"
		}
	`

	// Gets the saved weather
	testServer.Test("GET", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(getRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "{\"weather_report\":{\"vancouver\":{\"2020-04-18\":15}}}", responseBody)

	deleteRequestBody := `
		{
			"city": "vancouver"
		}
	`

	// Delete saved weather
	testServer.Test("DELETE", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(deleteRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "{\"message\":\"The weather was deleted succesfully!\"}", responseBody)

	// Try to get the saved weather and gets not found error
	testServer.Test("GET", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(getRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "{\"error\":\"Weather report not found\"}", responseBody)
	assert.NotContains(t, responseBody, "vancouver")
}
