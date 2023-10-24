#!/bin/bash

# Run producer and consumer apps using Go
echo "Starting producer and consumer apps using Go..."
cd ./consumer
go run main.go &
producer_pid=$!
echo "Producer PID: $producer_pid"
cd ../producer
go run main.go &
consumer_pid=$!
echo "Consumer PID: $consumer_pid"

# Wait for producer and consumer to start
echo "Waiting for producer and consumer to start..."
sleep 5

# Send a product request to the API
echo "Sending a product request..."
curl -X POST -H "Content-Type: application/json" -d '{
    "user_id": 1,
    "product_name": "Iphone",
    "product_description": "This iPhone is the latest model and packed with amazing features!",
    "product_images": [
        "https://raw.githubusercontent.com/krishnachowdaryvanam/products/main/Iphone13.jpg",
        "https://raw.githubusercontent.com/krishnachowdaryvanam/products/main/Iphone14.jpg",
        "https://raw.githubusercontent.com/krishnachowdaryvanam/products/main/Iphone15.jpg"
    ],
    "product_price": 60000
}' http://localhost:8080/products

# Wait for the consumer to consume all the messages
echo "Waiting for the consumer to consume all the messages..."
sleep 10

# Stop producer and consumer apps
echo "Stopping producer and consumer apps using Go..."
kill $producer_pid
kill $consumer_pid
echo "Stopped producer and consumer apps using Go."

# Optionally, stop the process listening to port 8080 (might not work on all systems)
echo "Stopping the process listening to port 8080..."
pkill -f "go run main.go"
echo "Stopped the process listening to port 8080."

