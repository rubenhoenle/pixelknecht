build:
	@go build -o pixelknecht ./cmd/pixelknecht/pixelknecht.go

test:
	@go test .

lint:
	@golangci-lint run
