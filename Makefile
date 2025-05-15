BINARY=gate-control
VERSION=1.0.0

all: build

build:
    CGO_ENABLED=0 go build -o ${BINARY} cmd/main.go

test:
    go test ./... -cover

lint:
    golangci-lint run

swagger:
    swag init

run:
    ./${BINARY}

clean:
    rm -f ${BINARY}

docker-build:
    docker build -t ${BINARY}:${VERSION} .

docker-run:
    docker run -p 8080:8080 ${BINARY}:${VERSION}

fmt:
    go fmt ./...

vet:
    go vet ./...