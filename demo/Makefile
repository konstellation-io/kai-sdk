SHELL=/bin/bash -e -o pipefail
COMPOSE_FILE=docker-compose.yml
DC=docker-compose -p sdk:local -f ${COMPOSE_FILE}
DOCKERFILE=exit-py/develop.Dockerfile

.DEFAULT_GOAL := help
include makefiles/help.mk
include makefiles/development.mk