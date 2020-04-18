package resources

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/felipecurvelo/weather-reporting-api/pkg/authorizer"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestAuthEndpoint_WithRightParams_ReturnOK(t *testing.T) {
	ctx := authorizer.NewContext(context.Background(), authorizer.NewAuthMock())
	testServer := api.NewTestServer(ctx, t).RegisterResource(&Auth{})

	requestBody := `
		{
			"name": "kirang",
			"password": "secret"
		}
	`

	testServer.Test("POST", "/auth/").WithBody(requestBody).Now()
	statusCode, responseBody := testServer.GetResponse()
	var actualModel authResponseModel
	err := json.Unmarshal([]byte(responseBody), &actualModel)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "M0CK3D_T0K3N", actualModel.Token)
}

func TestAuthEndpoint_WithInvalidName_ReturnError(t *testing.T) {
	testServer := api.NewTestServer(context.Background(), t).RegisterResource(&Auth{})

	requestBody := `
		{
			"name": "felipe",
			"password": "secret"
		}
	`

	testServer.Test("POST", "/auth/").WithBody(requestBody).Now()
	statusCode, responseBody := testServer.GetResponse()

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "{\"error\":\"Invalid Name\"}", responseBody)
}
