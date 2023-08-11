.DEFAULT_GOAL := help

# AutoDoc
# -------------------------------------------------------------------------
.PHONY: help
help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help

.PHONY: tidy
tidy: ## Run black, isort and codespell
	poetry --directory py-sdk run black py-sdk/src \
	&& poetry --directory py-sdk run isort py-sdk/src \
	&& poetry --directory py-sdk run codespell py-sdk/src -I py-sdk/dictionary.txt \
	--skip="*.git,*.json,kai_nats_msg_pb2.py,.venv,*.lock" \

.PHONY: protos
protos: ## Generate proto files
	protoc -I proto --python_out=py-sdk/src/sdk --go_out=go-sdk/protos --go_opt=paths=source_relative proto/kai_nats_msg.proto

.PHONY: generate_mocks
generate_mocks: ## Generate mocks
	cd go-sdk && go generate ./... && cd -

.PHONY: golint
golint: ## Run golint
	cd go-sdk && golangci-lint run ./... && cd -

.PHONY: gotest
gotest: ## Run tests
	cd go-sdk && go test ./... -cover -coverpkg=./... -coverprofile=coverage-unit.out --tags=unit

.PHONY: pytest
pytest: ## Run tests
	poetry --directory py-sdk run pytest py-sdk/src --cov --cov-report term-missing --cov-config=py-sdk/pyproject.toml --no-cov-on-fail
