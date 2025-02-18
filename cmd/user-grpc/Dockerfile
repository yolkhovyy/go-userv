FROM golang:1.23.6-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/user-grpc ./cmd/user-grpc
COPY internal ./internal
COPY contract ./contract
RUN go build -o main ./cmd/user-grpc

FROM alpine:3.20 AS release
WORKDIR /app
COPY --from=builder /app/main .
COPY cmd/user-grpc/config.yml .
EXPOSE 8080
CMD ["./main", "--config", "config.yml"]
