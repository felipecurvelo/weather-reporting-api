package main

import (
	"context"
	"fmt"

	"github.com/felipecurvelo/weather-reporting-api/pkg/auth"

	"github.com/felipecurvelo/weather-reporting-api/pkg/api"
	"github.com/felipecurvelo/weather-reporting-api/pkg/api/resources"
)

func main() {
	serverOptions := &api.ServerOptions{
		Port: 8080,
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{
		Name: "felipe",
	})

	server := api.NewServer(ctx, serverOptions)

	server.RegisterResource(&resources.Weather{}).
		Start()

	fmt.Printf("HTTP Server started on port: %d\n", serverOptions.Port)

	server.WaitForShutdownSignal().
		Close()

	fmt.Println("HTTP Server stopped")
}
