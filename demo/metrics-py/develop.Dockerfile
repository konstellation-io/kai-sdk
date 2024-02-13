FROM python:3.11.5-slim-bullseye

LABEL maintainer="Intelygenz - KAI Team"

ARG USER=kai
ARG UID=1001

ENV PATH="/root/.local/bin:$PATH" \
    POETRY_VIRTUALENVS_CREATE=false
ENV PYTHONUNBUFFERED=1

WORKDIR /tmp

COPY ["pyproject.toml", "poetry.lock", "./"] 

RUN apt update &&\
    apt install -yq --no-install-recommends curl  && \
    apt-get clean && apt-get autoremove -y && \
    useradd -m -b /sdk --shell /bin/bash --uid ${UID} ${USER} && \
    curl https://install.python-poetry.org -o poetry-install.py && \
    python poetry-install.py --version 1.5.1

RUN poetry install --only main --no-interaction --no-ansi

WORKDIR /app

RUN chown -R kai:0 /app \
    && chmod -R g+w /app \
    && mkdir /var/log/app -p \
    && chown -R kai:0 /var/log/app \
    && chmod -R g+w /var/log/app

USER ${USER}

COPY main.py process/main.py
COPY app.yaml app.yaml
COPY config.yaml config.yaml

CMD ["sh", "-c", "/usr/local/bin/python /app/process/main.py 2>&1 | tee -a /var/log/app/app.log"]