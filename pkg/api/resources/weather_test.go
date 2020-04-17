package resources

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestWeatherFirstEndpoint_ReturnWelcomeMessage(t *testing.T) {
	serverOptions := &api.ServerOptions{
		Port: 8080,
	}

	apiServer := api.NewServer(serverOptions).RegisterResource(&Weather{})
	httpServer := httptest.NewUnstartedServer(apiServer.GetHttpHandler())
	httpServer.Start()
	defer httpServer.Close()

	url := "http://" + httpServer.Listener.Addr().String()

	request, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	request.Header.Set("Accept", "application/json, text/plain, */*")

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	assert.NoError(t, err)

	responseBody, _ := ioutil.ReadAll(response.Body)

	assert.NoError(t, err)
	assert.NotNil(t, responseBody)
	assert.Equal(t, "Welcome! This is the first endpoint working!", string(responseBody))
}
