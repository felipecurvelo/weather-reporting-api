package resources

import (
	"context"
	"net/http"
	"testing"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestAuthEndpoint_WithRightParams_ReturnOK(t *testing.T) {
	testServer := api.NewTestServer(context.Background(), t).RegisterResource(&Auth{})

	requestBody := `
		{
			"ok": true
		}
	`

	testServer.Test("POST", "/auth/").WithBody(requestBody).Now()
	statusCode, responseBody := testServer.GetResponse()

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, requestBody, responseBody)
}
