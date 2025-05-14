BINARY=http-gate-control

build:
    go build -o ${BINARY} cmd/main.go

test:
    go test ./... -cover

lint:
    golangci-lint run

swagger:
    swag init

run: build
    ./${BINARY}