package main

import (
	"context"
	"fmt"

	"github.com/felipecurvelo/weather-reporting-api/pkg/authorizer"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/felipecurvelo/weather-reporting-api/pkg/api/resources"
)

func main() {
	serverOptions := &api.ServerOptions{
		Port: 8080,
	}

	ctx := context.Background()
	ctx = authorizer.NewContext(ctx, authorizer.NewAuth())

	server := api.NewServer(ctx, serverOptions)

	server.RegisterResource(&resources.Weather{}).
		RegisterResource(&resources.Auth{}).
		Start()

	fmt.Printf("HTTP Server started on port: %d\n", serverOptions.Port)

	server.WaitForShutdownSignal().
		Close()

	fmt.Println("HTTP Server stopped")
}
