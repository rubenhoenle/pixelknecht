build:
	@go build -o pixelknecht ./cmd/pixelknecht/pixelknecht.go

fmt:
	@gofmt -w .

test:
	@go test ./...

lint:
	@golangci-lint run
