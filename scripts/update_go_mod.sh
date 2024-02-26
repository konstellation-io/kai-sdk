#!/bin/bash

echo "Updating go-sdk"
cd go-sdk
rm go.sum
go mod tidy