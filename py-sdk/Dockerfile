FROM python:3.11-slim-bookworm

LABEL maintainer="Intelygenz - Konstellation Team"

ARG USER=kai

ENV PATH="/root/.local/bin:$PATH" \
    POETRY_VIRTUALENVS_CREATE=false

RUN apt update &&\
    apt install -yq --no-install-recommends curl  && \
    apt-get clean && apt-get autoremove -y && \
    useradd -m -b /sdk --shell /bin/bash --uid 1000 ${USER} && \
    curl https://install.python-poetry.org -o poetry-install.py && \
    python poetry-install.py --version 1.5.1

COPY ["pyproject.toml", "poetry.lock", "./"] 

RUN poetry install --only main --no-interaction --no-ansi

COPY --chown=1000 src /sdk/src

USER ${USER}