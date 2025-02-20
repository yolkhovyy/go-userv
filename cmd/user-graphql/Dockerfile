FROM golang:1.23.6-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/user-graphql ./cmd/user-graphql
COPY internal ./internal
COPY contract ./contract
RUN go build -o main ./cmd/user-graphql

FROM alpine:3.20 AS release
WORKDIR /app
COPY --from=builder /app/main .
COPY cmd/user-graphql/config.yml .
EXPOSE 8080
CMD ["./main", "--config", "config.yml"]
