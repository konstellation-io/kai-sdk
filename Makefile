generate_mocks:
	@echo "Generating mocks..."
	@cd go-sdk && go generate ./...
	@echo "Mocks generated"

generate_protobuf:
	@echo "Generating protobuf..."
	@protoc -I proto --go_out=go-sdk/protos --go_opt=paths=source_relative proto/kai_nats_msg.proto
	@echo "Protobuf generated"

build: generate_mocks generate_protobuf