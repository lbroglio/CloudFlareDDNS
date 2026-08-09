BINARY_NAME := cloudflareddns
BUILD_DIR=bin


.PHONY: test, build, clean

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/cloudflareddns/main.go

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...