BINARY_NAME := cloudflareddns
BUILD_DIR=bin


.PHONY: test, build, clean, build-static


build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/cloudflareddns/main.go

build-static:
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/cloudflareddns/main.go

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...