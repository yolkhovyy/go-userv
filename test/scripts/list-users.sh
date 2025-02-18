#!/bin/bash

# Params
page="${1:-1}"
limit="${2:-10}"
country="${3:-GB}"

url="http://localhost:8080/api/v1/users"

response=$(curl -s -G "$url" \
    --data-urlencode "page=$page" \
    --data-urlencode "limit=$limit" \
    --data-urlencode "country=$country")

echo "$response"
