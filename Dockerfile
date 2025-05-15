# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /gate-control cmd/main.go

# Final stage
FROM gcr.io/distroless/static-debian12

COPY --from=builder /gate-control /gate-control

EXPOSE 8080

CMD ["/gate-control"]