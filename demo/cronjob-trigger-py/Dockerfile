FROM python:3.11.5-slim-bullseye

LABEL maintainer="Intelygenz - Konstellation Team"

ARG USER=kai
ARG UID=1001

ENV PATH="/root/.local/bin:$PATH" \
    POETRY_VIRTUALENVS_CREATE=false

WORKDIR /tmp

COPY ["pyproject.toml", "poetry.lock", "./"] 

RUN apt update &&\
    apt install -yq --no-install-recommends curl  && \
    apt-get clean && apt-get autoremove -y && \
    useradd -m -b /sdk --shell /bin/bash --uid ${UID} ${USER} && \
    curl https://install.python-poetry.org -o poetry-install.py && \
    python poetry-install.py --version 1.5.1

RUN poetry install --only main --no-interaction --no-ansi

USER ${USER}

WORKDIR /app

COPY main.py process/main.py
COPY app.yaml app.yaml
COPY config.yaml config.yaml

CMD ["python","/app/process/main.py"]
