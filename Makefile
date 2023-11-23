.DEFAULT_GOAL := help

# AutoDoc
# -------------------------------------------------------------------------
.PHONY: help
help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help

.PHONY: protos
protos: ## Generate proto files
	protoc -I proto --python_out=py-sdk/sdk/sdk --mypy_out=py-sdk/sdk/sdk --go_out=go-sdk/protos --go_opt=paths=source_relative proto/kai_nats_msg.proto

.PHONY: generate_mocks
generate_mocks: ## Generate mocks
	cd go-sdk && go generate ./... && cd -

.PHONY: gotidy
gotidy: ## Run golangci-lint, goimports and gofmt
	cd go-sdk && golangci-lint run ./... && goimports -w  . && gofmt -s -w -e -d . && cd -

.PHONY: gotest
gotest: ## Run integration and unit tests
	cd go-sdk && go test ./... -cover -coverpkg=./... --tags=unit,integration

.PHONY: pytest
pytest: ## Run unit tests
	cd py-sdk && poetry run pytest sdk runner --cov --cov-report=term-missing

.PHONY: pytidy
pytidy: ## Run black, isort and codespell
	poetry --directory py-sdk run black py-sdk \
	&& poetry --directory py-sdk run isort py-sdk \
	&& poetry --directory py-sdk run codespell py-sdk -I py-sdk/dictionary.txt \
	--skip="*.git,*.json,kai_nats_msg_pb2.py,.venv,*.lock,__init__.py" \
	&& poetry --directory py-sdk run black demo \
	&& poetry --directory py-sdk run isort demo \
	&& poetry --directory py-sdk run codespell demo -I py-sdk/dictionary.txt \
	--skip="*.git,*.json,kai_nats_msg_pb2.py,.venv,*.lock,__init__.py" \

.PHONY: mypy
mypy: ## Run mypy
	poetry --directory py-sdk run mypy --pretty --warn-redundant-casts --warn-unused-ignores --warn-unreachable --disallow-untyped-decorators --disallow-incomplete-defs --disallow-untyped-calls --check-untyped-defs --disallow-incomplete-defs --python-version 3.11 py-sdk --config-file py-sdk/pyproject.toml
