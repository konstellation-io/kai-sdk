#!/bin/bash

SOURCE_PATH="$PWD/src"

protoc -I=. \
  --python_out="$SOURCE_PATH/training-py/proto" \
  --mypy_out="$SOURCE_PATH/training-py/proto" \
  --python_out="$SOURCE_PATH/splitter/proto" \
  --mypy_out="$SOURCE_PATH/splitter/proto" \
  *.proto

protoc -I=. \
  --go_out="$SOURCE_PATH/training-go/proto" \
  --go_out="$SOURCE_PATH/validation/proto" \
  --go_out="$SOURCE_PATH/exit/proto" \
  --go_opt=paths=source_relative *.proto \

echo "Done"
