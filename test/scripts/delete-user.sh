#!/bin/bash

# User ID
uuid="${1:-}"

url="http://localhost:8080/api/v1/user/$uuid"

response=$(curl -s -X DELETE "$url")

echo "$response"
