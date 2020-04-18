package api

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestServer struct {
	apiServer    *Server
	httpServer   *httptest.Server
	httpRequest  *http.Request
	httpResponse *http.Response
	errors       []error
	t            *testing.T
}

func (ts *TestServer) RegisterResource(resource Resource) *TestServer {
	ts.apiServer.RegisterResource(resource)
	return ts
}

func (ts *TestServer) Test(method string, endpointURL string) *TestServer {
	ts.httpServer = httptest.NewUnstartedServer(ts.apiServer.GetHttpHandler())

	url := "http://" + ts.httpServer.Listener.Addr().String() + endpointURL

	request, err := http.NewRequest(method, url, nil)
	assert.NoError(ts.t, err)

	request.Header.Set("Accept", "application/json, text/plain, */*")

	ts.httpRequest = request
	return ts
}

func (ts *TestServer) WithHeader(key string, value string) *TestServer {
	ts.httpRequest.Header.Set(key, value)
	return ts
}

func (ts *TestServer) WithBody(requestBody string) *TestServer {
	ts.httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(requestBody)))
	ts.httpRequest.ContentLength = int64(len([]byte(requestBody)))
	return ts
}

func (ts *TestServer) Now() *TestServer {
	ts.httpServer.Start()
	defer ts.httpServer.Close()

	httpClient := &http.Client{}
	response, err := httpClient.Do(ts.httpRequest)
	assert.NoError(ts.t, err)

	ts.httpResponse = response
	return ts
}

func (ts *TestServer) GetResponse() (int, string) {
	responseBody, err := ioutil.ReadAll(ts.httpResponse.Body)
	assert.NoError(ts.t, err)
	return ts.httpResponse.StatusCode, string(responseBody)
}

func NewTestServer(ctx context.Context, t *testing.T) *TestServer {
	serverOptions := &ServerOptions{
		Port: 8080,
	}

	apiServer := NewServer(ctx, serverOptions)

	return &TestServer{
		apiServer: apiServer,
		t:         t,
	}
}
