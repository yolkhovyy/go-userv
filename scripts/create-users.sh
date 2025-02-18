#!/bin/bash

generate_json_payload() {
    local i="$1"
    local country="$2"
    jq -n --arg i "$i" --arg country "$country" '{
        firstName: "John",
        lastName: ("Doe, " + $i),
        nickname: ("johnd" + $i),
        password: "securepassword",
        email: ("john.doe." + $i + "@example.com"),
        country: $country
    }'
}

num_users="${1:-42}"
country="${2:-GB}"

for ((i = 1; i <= num_users; i++)); do
    {
        json_payload=$(generate_json_payload "$i" "$country")

        response=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/api/v1/user \
            -H "Content-Type: application/json" \
            -d "$json_payload")

        if [[ "$response" -eq 0 ]]; then
            echo "Curl request failed, server might be down"
            exit 1
        elif [[ "$response" -eq 201 ]]; then
            echo "User #$i created successfully"
        else
            echo "Failed to create user #$i, HTTP status: $response"
        fi
    } &
done

wait
