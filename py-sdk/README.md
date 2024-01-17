# AIO Python SDK

AIO SDK's implementation in Python.


## How it works

### SDK

The SDK can be used in a AIO workflows for working with Python code

Once the Python SDK is deployed, it connects to NATS and it subscribes permanently to an input subject. Each node knows to which subject it has to subscribe and also to which subject it has to send messages, since the K8s manager tells it with environment variables. It is important to note that the nodes use a queue subscription, which allows load balancing of messages when there are multiple replicas when using Task and Exit runners but not with Trigger runners

When a new message is published in the input subject of a node, it passes it down to a handler function, along with a context object formed by variables and useful methods for processing data. This handler is the solution implemented by the client and given in the krt file generated. Once executed, the result will be taken and transformed into a NATS message that will then be published to the next node's subject (indicated by an environment variable). After that, the node ACKs the message manually

## Development

- Build the proto files with `make protos`

- Install the dependencies with `poetry install --group dev`

If you don't have poetry installed (you must have python 3.11 installed in your system):

`python3 -m pip install --user poetry`

## Tests

Just run `make pytest` from the root folder


## Proto

### Python

- `wget https://github.com/protocolbuffers/protobuf/releases/download/v23.4/protoc-23.4-linux-x86_64.zip`
- `unzip -o protoc-23.4-linux-x86_64.zip -d /usr/local bin/protoc`
- `unzip -o protoc-23.4-linux-x86_64.zip -d /usr/local 'include/*'`

### Python SDK/PY-SDK package release

- First go to SDK or PY-SDK folder and run poetry build
- Then run `poetry publish --repository my-gitlab -u DEPLOY_TOKEN_USER -p DEPLOY_TOKEN`
- Deploy token can be generated in `https://gitlab.zenith.igzdev.com/zenith-platform/mlops-aiorchestrator/aio-sdk/-/settings/repository`
 deploy tokens section, ask your admin for the one used in this repository
- You can find the package in `https://gitlab.zenith.igzdev.com/zenith-platform/mlops-aiorchestrator/aio-sdk/-/packages`

### Using Python SDK/PY-SDK in a project

- The dependency is added like `py-sdk = {source = "aio-sdk", version = "1.0.0"}`
- Then we need to add:
```
[[tool.poetry.source]]
name = "aio-sdk"
url = "https://gitlab.zenith.igzdev.com/api/v4/projects/498/packages/pypi/simple"
priority = "supplemental"
```
- Poetry will use the credentials stored in `~/.netrc` that needs to be defined as:
```
machine gitlab.zenith.igzdev.com login GITLAB_USER password GITLAB_USER_TOKEN
```
or you need to provide them in a different way
- Gitlab user token can be obtained in `https://gitlab.zenith.igzdev.com/-/user_settings/personal_access_tokens`