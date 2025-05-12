# Билд-стадия
FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /gate-service

# Финальная стадия
FROM gcr.io/distroless/static-debian12
COPY --from=builder /gate-service /gate-service
COPY docs /docs

EXPOSE 8080

ENTRYPOINT ["/gate-service"]