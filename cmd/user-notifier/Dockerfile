FROM golang:1.23.6-alpine3.20 AS builder
ARG APP_VERSION
ARG LDFLAGS="-X github.com/yolkhovyy/go-userv/cmd/user-notifier/version.Tag=${APP_VERSION}"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/user-notifier ./cmd/user-notifier
COPY internal ./internal
RUN go build -ldflags="${LDFLAGS}" -o main ./cmd/user-notifier && \
    go test -tags test -coverpkg=./... -c -o main.test ./cmd/user-notifier -args --config config.yml

FROM alpine:3.20 AS release
WORKDIR /app
COPY --from=builder /app/main .
COPY cmd/user-notifier/config.yml .
CMD ["./main", "--config", "config.yml"]

FROM alpine:3.20 AS test
WORKDIR /app
COPY --from=builder /app/main.test .
COPY cmd/user-notifier/config.yml .
EXPOSE 8080
CMD ["./main.test", "-test.run", "^TestRunMain$", "-test.coverprofile", "user-notifier.cov", "-test.v"]
