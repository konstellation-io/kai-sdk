# KAI SDK

This repository contains all available languages SDKs and Runners for the [KAI](https://github.com/konstellation-io/kai).

## Runners coverage

|      Component       |                       Coverage                       |                       Bugs                       |               Maintainability Rating               |                      Go report                     |
|:--------------------:|:----------------------------------------------------:|:------------------------------------------------:| :------------------------------------------------: | :------------------------------------------------: |
|        GO SDK        | [![coverage][go-sdk-coverage]][go-sdk-coverage-link] | [![bugs][go-sdk-bugs]][go-sdk-bugs-link] | [-![mr][go-sdk-mr]][go-sdk-mr-link] | - |
|        Py SDK        | [![coverage][py-sdk-coverage]][py-sdk-coverage-link] | [![bugs][py-sdk-bugs]][py-sdk-bugs-link] | [-![mr][py-sdk-mr]][py-sdk-mr-link] | - |

[go-sdk-coverage]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_go-sdk&metric=coverage
[go-sdk-coverage-link]: https://sonarcloud.io/dashboard?id=konstellation-io_go-sdk
[go-sdk-bugs]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_go-sdk&metric=bugs
[go-sdk-bugs-link]: https://sonarcloud.io/dashboard?id=konstellation-io_go-sdk
[go-sdk-loc]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_go-sdk&metric=ncloc
[go-sdk-loc-link]: https://sonarcloud.io/dashboard?id=konstellation-io_go-sdk
[go-sdk-mr]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_go-sdk&metric=sqale_rating
[go-sdk-mr-link]: https://sonarcloud.io/dashboard?id=konstellation-io_go-sdk
[py-sdk-coverage]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_py-sdk&metric=coverage
[py-sdk-coverage-link]: https://sonarcloud.io/dashboard?id=konstellation-io_py-sdk
[py-sdk-bugs]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_py-sdk&metric=bugs
[py-sdk-bugs-link]: https://sonarcloud.io/dashboard?id=konstellation-io_py-sdk
[py-sdk-loc]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_py-sdk&metric=ncloc
[py-sdk-loc-link]: https://sonarcloud.io/dashboard?id=konstellation-io_py-sdk
[py-sdk-mr]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_py-sdk&metric=sqale_rating
[py-sdk-mr-link]: https://sonarcloud.io/dashboard?id=konstellation-io_py-sdk

## Protobuf

All components receive and send a `KaiNatsMessage` protobuf.
To generate the protobuf code, the `protoc` compiler must be installed.

### Installation

- `wget https://github.com/protocolbuffers/protobuf/releases/download/v23.4/protoc-23.4-linux-x86_64.zip`
- `unzip -o protoc-23.4-linux-x86_64.zip -d /usr/local bin/protoc`
- `unzip -o protoc-23.4-linux-x86_64.zip -d /usr/local 'include/*'`

### Regenerating protos

```
make protos
```
