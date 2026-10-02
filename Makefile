.PHONY: keys tidy test build server token balance transfer clean

# Generate the RSA keypair into ./keys
keys:
	go run ./cmd/snapkeys -out keys

# Download/refresh module dependencies.
tidy:
	go mod tidy

test:
	go test ./...

build:
	go build ./...

# Run the mock SNAP server (foreground).
server:
	go run ./cmd/mockserver

# Client commands (server must be running in another terminal).
token:
	go run ./cmd/snapclient token

balance:
	go run ./cmd/snapclient balance 115471119

transfer:
	go run ./cmd/snapclient transfer 115471119 7382382957893840 10000.00

clean:
	rm -rf keys bin