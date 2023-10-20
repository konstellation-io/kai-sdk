# KRE Runners

This repo contains all language runners for [KRE](https://github.com/konstellation-io/kre).

## Runners coverage

|      Component       |                       Coverage                       |                       Bugs                       |               Maintainability Rating               |                      Go report                     |
|:--------------------:|:----------------------------------------------------:|:------------------------------------------------:| :------------------------------------------------: | :------------------------------------------------: |
|        GO SDK        | [![coverage][go-sdk-coverage]][go-sdk-coverage-link] | [![bugs][go-sdk-bugs]][go-sdk-bugs-link] | [-![mr][go-sdk-mr]][go-sdk-mr-link] | - |
| KRT Files Downloader | [![coverage][krt-fd-coverage]][krt-fd-coverage-link] |     [![bugs][krt-fd-bugs]][krt-fd-bugs-link]     |         [![mr][krt-fd-mr]][krt-fd-mr-link]         | - |


[go-sdk-coverage]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_go-sdk&metric=coverage
[go-sdk-coverage-link]: https://sonarcloud.io/dashboard?id=konstellation-io_go-sdk
[go-sdk-bugs]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_kre_py&metric=bugs
[go-sdk-bugs-link]: https://sonarcloud.io/dashboard?id=konstellation-io_go-sdk
[go-sdk-loc]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_go-sdk&metric=ncloc
[go-sdk-loc-link]: https://sonarcloud.io/dashboard?id=konstellation-io_go-sdk
[go-sdk-mr]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_go-sdk&metric=sqale_rating
[go-sdk-mr-link]: https://sonarcloud.io/dashboard?id=konstellation-io_go-sdk
[krt-fd-coverage]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_krt_files_downloader&metric=coverage
[krt-fd-coverage-link]: https://sonarcloud.io/dashboard?id=konstellation-io_krt_files_downloader
[krt-fd-bugs]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_krt_files_downloader&metric=bugs
[krt-fd-bugs-link]: https://sonarcloud.io/dashboard?id=konstellation-io_krt_files_downloader
[krt-fd-loc]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_krt_files_downloader&metric=ncloc
[krt-fd-loc-link]: https://sonarcloud.io/dashboard?id=konstellation-io_krt_files_downloader
[krt-fd-mr]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_krt_files_downloader&metric=sqale_rating
[krt-fd-mr-link]: https://sonarcloud.io/dashboard?id=konstellation-io_krt_files_downloader

## Protobuf

All components receive and send a `KaiNatsMessage` protobuf.
To generate the protobuf code, the `protoc` compiler must be installed.
Use the following command to generate the code:

```
make build
```
