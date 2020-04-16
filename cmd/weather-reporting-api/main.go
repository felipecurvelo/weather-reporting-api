package main

import (
	"fmt"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/felipecurvelo/weather-reporting-api/pkg/api/resources"
)

func main() {
	serverOptions := &api.ServerOptions{
		Port: 8080,
	}

	server := api.NewServer(serverOptions)

	server.RegisterResource(&resources.Weather{}).
		Start()

	fmt.Printf("HTTP Server started on port: %d\n", serverOptions.Port)

	server.WaitForShutdownSignal().
		Close()

	fmt.Println("HTTP Server stopped")
}
