BINARY_NAME=klotski-go

all: test build

test:
	go test -v ./pkg

clean:
	rm -rf build/$(BINARY_NAME)

build: clean
	go build -ldflags="-s -w" -o build/$(BINARY_NAME) cmd/main.go
	
run-cli:
	./build/$(BINARY_NAME) -mode cli

run-http:
	./build/$(BINARY_NAME) -mode http