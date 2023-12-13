#!/bin/bash

echo "Updating go-sdk"
cd go-sdk
rm go.sum
go mod tidy

echo "Updating demo"

echo "Updating cronjob-trigger"
cd ../demo/cronjob-trigger
rm go.sum
go mod tidy

echo "Updating exit"
cd ../exit
rm go.sum
go mod tidy

echo "Updating nats-trigger"
cd ../nats-trigger
rm go.sum
go mod tidy

echo "Updating rest-trigger"
cd ../rest-trigger
rm go.sum
go mod tidy

echo "Updating metrics"
cd ../metrics
rm go.sum
go mod tidy

echo "Updating task"
cd ../task
rm go.sum
go mod tidy