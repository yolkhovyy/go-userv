FROM golang:1.23.6-alpine3.20 AS builder
ARG APP_VERSION
ARG LDFLAGS="-X github.com/yolkhovyy/go-userv/cmd/user-graphql/version.Tag=${APP_VERSION}"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/user-graphql ./cmd/user-graphql
COPY internal ./internal
COPY contract ./contract
RUN go build -ldflags="${LDFLAGS}" -o main ./cmd/user-graphql && \
    go test -tags test -coverpkg=./... -c -o main.test ./cmd/user-graphql -args --config config.yml

FROM alpine:3.20 AS release
WORKDIR /app
COPY --from=builder /app/main .
COPY cmd/user-graphql/config.yml .
EXPOSE 8080
CMD ["./main", "--config", "config.yml"]

FROM alpine:3.20 AS test
WORKDIR /app
COPY --from=builder /app/main.test .
COPY cmd/user-graphql/config.yml .
EXPOSE 8080
CMD ["./main.test", "-test.run", "^TestRunMain$", "-test.coverprofile", "user-graphql.cov", "-test.v"]

