#!/bin/bash

echo "Updating py-sdk"
cd py-sdk/sdk
poetry update
cd ..
poetry update