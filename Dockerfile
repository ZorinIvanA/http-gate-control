# Stage 1: Build
FROM golang:1.24 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /http-gate-control cmd/main.go

# Stage 2: Run
FROM gcr.io/distroless/static-debian12
COPY --from=builder /http-gate-control /http-gate-control
EXPOSE 8080
CMD ["/http-gate-control"]