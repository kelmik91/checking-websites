export GO111MODULE=on
export GOOS=linux
export GOARCH=amd64

TARGET ?= checkerSite

#golangci-lint:
#	find -type f -name "*.go" | grep -v '.*\.pb\.go' | grep -v '\/[0-9a-z_]*.go' && echo "Files should be named in snake case" && exit 1 || echo "All files named in snake case"
#	golangci-lint version
#	golangci-lint run

buildMainChecker:
	go build -o bin/$(TARGET) ./cmd/main

buildSslChecker:
	go build -o bin/checkerSsl ./cmd/ssl

buildWeather:
	go build -o bin/weather ./cmd/weather

tidy:
	go mod tidy

clean:
	rm -rf bin/$(TARGET)

download:
	go mod download