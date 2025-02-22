#!/bin/bash

SERVICES=("user-graphql" "user-grpc" "user-notifier" "user-rest")

echo "⚙️ Running integration tests"
go test -v -count=1 --tags=integration_tests ./test/integration/...

echo "⚙️ Exporting coverage profiles"
mkdir -p coverage

for SERVICE in "${SERVICES[@]}"; do
    CONTAINER=$(docker compose ps -q $SERVICE)

    echo "⚙️ Stopping $SERVICE $CONTAINER"
    docker compose stop --timeout 30 ${SERVICE}

    echo "⚙️ Exporting coverage $SERVICE $CONTAINER"
    docker export ${CONTAINER} > ${SERVICE}.tar
    tar -xf ${SERVICE}.tar app/${SERVICE}.cov
    cat app/${SERVICE}.cov | grep -v -e "mock_" -v -e "test" > coverage/${SERVICE}.cov
    
    rm ${SERVICE}.tar
done

echo "🗑 Cleaning up"
rm -rf app/
