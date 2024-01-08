# Testing manually with a docker image

- Inside each metrics demo process folder you can find two files:
    - otel-collector-config.yaml
    - otelCollectorDockerfile

Rename otelCollectorDockerfile to Dockerfile and then you just need to get inside that folder and run `docker build -t test-otl-collector .` and then `docker run -p 4317:4317 test-otl-collector`