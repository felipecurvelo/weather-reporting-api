build:
	go build ./cmd/weather-reporting-api/

run: build
	./weather-reporting-api

test:
	go clean -testcache
	go test -v -parallel 5 ./...