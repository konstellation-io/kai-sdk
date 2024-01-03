.PHONY: build
build: ##@development Build containers. Needed only after changes in poetry dependencies.
build: args?= -f ${DOCKERFILE}
build:
	DOCKER_BUILDKIT=1 docker build -t "sdk:local" \
		--progress plain \
		${args} .

.PHONY: clean
clean: ##@development  Stop and remove containers created by up command
	${DC} down --remove-orphans

.PHONY: logs
logs: ##@development Show logs for the current project
	${DC} logs -f

.PHONY: server
server: ##@development Start a SDK instance with NATS and minio
server:
	${DC} up -d app

.PHONY: shell
shell: ##@development Starts a bash shell
	${DC} run --rm --entrypoint="/bin/bash" server

.PHONY: status
status: ##@development List images used by the created container, list containers and display the running processes.
	@echo
	${DC} images
	@echo
	${DC} ps
	@echo
	${DC} top