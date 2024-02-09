#!/bin/bash

echo "Updating py-sdk"
cd py-sdk/sdk
poetry update
cd ..
poetry update

echo "Updating demo"

echo "Updating cronjob-trigger-py"
cd ../demo/cronjob-trigger-py
poetry update

echo "Updating exit-py"
cd ../exit-py
poetry update

echo "Updating nats-trigger-py"
cd ../nats-trigger-py
poetry update

echo "Updating metrics-py"
cd ../metrics-py
poetry update

echo "Updating task-py"
cd ../task-py
poetry update

echo "Updating redis-py"
cd ../redis-py
poetry update