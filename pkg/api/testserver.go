package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type TestServer struct {
	apiServer *Server
}

func (ts *TestServer) RegisterResource(resource Resource) *TestServer {
	ts.apiServer.RegisterResource(resource)
	return ts
}

func (ts *TestServer) CallEndpoint(method string, endpointURL string) (string, int, error) {
	httpServer := httptest.NewUnstartedServer(ts.apiServer.GetHttpHandler())
	httpServer.Start()
	defer httpServer.Close()

	url := "http://" + httpServer.Listener.Addr().String() + endpointURL

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", 0, err
	}

	request.Header.Set("Accept", "application/json, text/plain, */*")

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)

	responseBody, _ := ioutil.ReadAll(response.Body)

	return string(responseBody), response.StatusCode, nil
}

func NewTestServer() *TestServer {
	serverOptions := &ServerOptions{
		Port: 8080,
	}

	apiServer := NewServer(serverOptions)

	return &TestServer{
		apiServer: apiServer,
	}
}
