FROM golang:1.23.6-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/user-rest ./cmd/user-rest
COPY internal ./internal
COPY contract/dto ./contract/dto
RUN go build -o main ./cmd/user-rest

FROM alpine:3.20 AS release
WORKDIR /app
COPY --from=builder /app/main .
COPY cmd/user-rest/config.yml .
EXPOSE 8080
CMD ["./main", "--config", "config.yml"]
