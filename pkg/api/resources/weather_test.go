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

func TestWeatherSave_WithEmptyAuthHeader_ReturnUnauthorized(t *testing.T) {
	ctx := authorizer.NewContext(context.Background(), authorizer.NewAuthMock())
	ctx = weathermanager.NewContext(ctx, weathermanager.New())

	testServer := api.NewTestServer(ctx, t).
		RegisterResource(&Weather{})

	testServer.Test("POST", "/weather/").Now()
	statusCode, responseBody := testServer.GetResponse()

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "{\"error\":\"Error validating auth token (Empty Token)\"}", responseBody)
}

func TestWeatherSave_WithInvalidAuthToken_ReturnUnauthorized(t *testing.T) {
	ctx := authorizer.NewContext(context.Background(), authorizer.NewAuthMock())
	ctx = weathermanager.NewContext(ctx, weathermanager.New())

	testServer := api.NewTestServer(ctx, t).
		RegisterResource(&Weather{})

	testServer.Test("POST", "/weather/").
		WithHeader("Authorization", "invalid_token").
		Now()
	statusCode, responseBody := testServer.GetResponse()

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "{\"error\":\"Error validating auth token (Invalid Token)\"}", responseBody)
}

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

func TestWeatherSave_WithInvalidDate_ReturnError(t *testing.T) {
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
				"date": "invalid_date",
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

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Contains(t, responseBody, "Invalid date")
}

func TestWeatherSave_WithInvalidDateRange_ReturnError(t *testing.T) {
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
				"date": "2020-02-31",
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

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Contains(t, responseBody, "Invalid date")
}

func TestWeatherGet_WithEmptyAuthHeader_ReturnUnauthorized(t *testing.T) {
	ctx := authorizer.NewContext(context.Background(), authorizer.NewAuthMock())
	ctx = weathermanager.NewContext(ctx, weathermanager.New())

	testServer := api.NewTestServer(ctx, t).
		RegisterResource(&Weather{})

	testServer.Test("GET", "/weather/").Now()
	statusCode, responseBody := testServer.GetResponse()

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "{\"error\":\"Error validating auth token (Empty Token)\"}", responseBody)
}

func TestWeatherGet_WithInvalidAuthToken_ReturnUnauthorized(t *testing.T) {
	ctx := authorizer.NewContext(context.Background(), authorizer.NewAuthMock())
	ctx = weathermanager.NewContext(ctx, weathermanager.New())

	testServer := api.NewTestServer(ctx, t).
		RegisterResource(&Weather{})

	testServer.Test("GET", "/weather/").
		WithHeader("Authorization", "invalid_token").
		Now()
	statusCode, responseBody := testServer.GetResponse()

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "{\"error\":\"Error validating auth token (Invalid Token)\"}", responseBody)
}

func TestWeatherGet_ReturnOK(t *testing.T) {
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

func TestWeatherGet_WithEmptyCity_ReturnError(t *testing.T) {
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

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Contains(t, responseBody, "Empty city")
}

func TestWeatherGet_WithEmptyInitialDate_ReturnError(t *testing.T) {
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
			"end_date": "2020-04-30"
		}
	`

	// Gets the saved weather
	testServer.Test("GET", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(getRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Contains(t, responseBody, "Empty initial date")
}

func TestWeatherGet_WithInvalidInitialDate_ReturnOK(t *testing.T) {
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
			"initial_date": "invalid_date",
			"end_date": "2020-04-30"
		}
	`

	// Gets the saved weather
	testServer.Test("GET", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(getRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Contains(t, responseBody, "Invalid initial date")
}

func TestWeatherGet_WithInvalidInitialDateRange_ReturnOK(t *testing.T) {
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
			"initial_date": "2020-02-31",
			"end_date": "2020-04-30"
		}
	`

	// Gets the saved weather
	testServer.Test("GET", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(getRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Contains(t, responseBody, "Invalid initial date")
}

func TestWeatherGet_WithEmptyEndDate_ReturnOK(t *testing.T) {
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
			"initial_date": "2020-04-01"
		}
	`

	// Gets the saved weather
	testServer.Test("GET", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(getRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Contains(t, responseBody, "Empty end date")
}

func TestWeatherGet_WithInvalidEndDate_ReturnOK(t *testing.T) {
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
			"end_date": "invalid_date"
		}
	`

	// Gets the saved weather
	testServer.Test("GET", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(getRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Contains(t, responseBody, "Invalid end date")
}

func TestWeatherGet_WithInvalidEndDateRange_ReturnOK(t *testing.T) {
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
			"end_date": "2020-02-31"
		}
	`

	// Gets the saved weather
	testServer.Test("GET", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(getRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Contains(t, responseBody, "Invalid end date")
}

func TestWeatherGet_WithInvalidDateRange_ReturnOK(t *testing.T) {
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
			"initial_date": "2020-04-20",
			"end_date": "2020-04-01"
		}
	`

	// Gets the saved weather
	testServer.Test("GET", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(getRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Contains(t, responseBody, "Invalid date range")
}

func TestWeatherDelete_WithEmptyAuthHeader_ReturnUnauthorized(t *testing.T) {
	ctx := authorizer.NewContext(context.Background(), authorizer.NewAuthMock())
	ctx = weathermanager.NewContext(ctx, weathermanager.New())

	testServer := api.NewTestServer(ctx, t).
		RegisterResource(&Weather{})

	testServer.Test("DELETE", "/weather/").Now()
	statusCode, responseBody := testServer.GetResponse()

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "{\"error\":\"Error validating auth token (Empty Token)\"}", responseBody)
}

func TestWeatherDelete_WithInvalidAuthToken_ReturnUnauthorized(t *testing.T) {
	ctx := authorizer.NewContext(context.Background(), authorizer.NewAuthMock())
	ctx = weathermanager.NewContext(ctx, weathermanager.New())

	testServer := api.NewTestServer(ctx, t).
		RegisterResource(&Weather{})

	testServer.Test("DELETE", "/weather/").
		WithHeader("Authorization", "invalid_token").
		Now()
	statusCode, responseBody := testServer.GetResponse()

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "{\"error\":\"Error validating auth token (Invalid Token)\"}", responseBody)
}

func TestWeatherDelete_ReturnEmpty(t *testing.T) {
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

	// Gets the saved weather to verify
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

	// Try to get the saved weather and gets error
	testServer.Test("GET", "/weather/").
		WithHeader("Authorization", tokenObj["token"].(string)).
		WithBody(getRequestBody).
		Now()
	statusCode, responseBody = testServer.GetResponse()

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "{\"error\":\"Weather report not found\"}", responseBody)
	assert.NotContains(t, responseBody, "vancouver")
}
