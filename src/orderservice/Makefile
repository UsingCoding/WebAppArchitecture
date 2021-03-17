export APP_CMD_NAME = orderservice

all: build

build:
	bin/go-build.sh "cmd" "bin/$(APP_CMD_NAME)" $(APP_CMD_NAME)

run:
	./bin/${APP_CMD_NAME}

test:
	go test ./...

check:
	golangci-lint run