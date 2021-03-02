export APP_CMD_NAME = orderservice

all: build run

build:
	go build -o ./bin/${APP_CMD_NAME} ./cmd/${APP_CMD_NAME}/main.go

run:
	./bin/${APP_CMD_NAME}

test:
	go test ./...
