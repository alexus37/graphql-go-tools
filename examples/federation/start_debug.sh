#!/bin/bash

# function cleanup {
#     kill "$ACCOUNTS_PID"
#     kill "$PRODUCTS_PID"
#     kill "$REVIEWS_PID"
# }
# trap cleanup EXIT
kill -9 $(lsof -t -i:4001)
kill -9 $(lsof -t -i:4002)
kill -9 $(lsof -t -i:4003)

echo "Building services"
go build -o /tmp/srv-accounts ./accounts
go build -o /tmp/srv-products ./products
go build -o /tmp/srv-reviews ./reviews
go build -o /tmp/srv-gateway ./gateway

echo "Starting services"
/tmp/srv-accounts &
ACCOUNTS_PID=$!

/tmp/srv-products &
PRODUCTS_PID=$!

/tmp/srv-reviews &
REVIEWS_PID=$!

echo "Please start vscode debugger now"

