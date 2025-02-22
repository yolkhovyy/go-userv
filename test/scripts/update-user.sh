#!/bin/bash

generate_json_payload() {
    local country="$1"
    jq -n --arg country "$country" '{
        "firstName": "John",
        "lastName": "Doe, 424242",
        "nickname": "johnd424242",
        "password": "securepassword",
        "email": "john.doe.424242@example.com",
        "country": $country
    }'
}

userID="${1}"
country="${2:-AU}"

if [[ -z "$userID" ]]; then
    echo "Error: No user ID provided."
    echo "Usage: $0 <userID> [country]"
    exit 1
fi

json_payload=$(generate_json_payload "$country")

response=$(curl -s -o /dev/null -w "%{http_code}" -X PUT "http://localhost:8080/api/v1/user/$userID" \
    -H "Content-Type: application/json" \
    -d "$json_payload")

if [[ -z "$response" || "$response" -eq 0 ]]; then
    echo "Curl request failed (network error or server down)"
    exit 1
elif [[ "$response" -eq 200 || "$response" -eq 202 ]]; then
    echo "User #$userID updated successfully"
else
    echo "Failed to update user #$userID, HTTP status: $response"
fi
