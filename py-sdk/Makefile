.DEFAULT_GOAL := help

# AutoDoc
# -------------------------------------------------------------------------
.PHONY: help
help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help


.PHONY: tidy
tidy: ## Run black, isort and codespell
	poetry run black src \
	&& poetry run isort src \
	&& poetry run codespell src -I dictionary.txt \
	--skip="*.git,*.json,public_input_pb2.py,public_input_pb2_grpc.py,.venv,*.lock" \

.PHONY: protos
protos: ## Generate proto files
	poetry -- run python -m grpc_tools.protoc -I="./protos" \
	--grpc_python_out="./src" \
	./protos/public_input.proto \
	&& poetry run python -m grpc_tools.protoc -I="./protos" \
	--grpc_python_out="./src" \
	./protos/internal_nodes.proto

.PHONY: test
test: ## Run tests
	poetry run pytest src/tests/src/ --cov --cov-report term-missing --cov-config=pyproject.toml

.PHONY: docker
docker: ## Build and Run py-sdk docker
	docker build . -t py-sdk:latest && docker run --name py-sdk -ti py-sdk:latest 

.PHONY: attach-docker
attach-docker: ## Attach ssh to py-sdk docker
	docker exec -it dni-extractor sh

.PHONY: clean-all-docker
clean-all-docker: ## Clean all docker
	docker rm $$(docker ps -a -q) -f && docker rmi $$(docker images -a -q) -f