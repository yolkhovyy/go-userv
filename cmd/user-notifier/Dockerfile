FROM golang:1.23.6-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/user-notifier ./cmd/user-notifier
COPY internal ./internal
RUN go build -o main ./cmd/user-notifier

FROM alpine:3.20 AS release
WORKDIR /app
COPY --from=builder /app/main .
COPY cmd/user-notifier/config.yml .
CMD ["./main", "--config", "config.yml"]
