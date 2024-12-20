# Stage 1: Build the binary
FROM golang:1.23 AS builder

WORKDIR /app

COPY .. .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o client ./cmd/client

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/client /app/client
COPY --from=builder /app/config /app/config

CMD ["/app/client"]
