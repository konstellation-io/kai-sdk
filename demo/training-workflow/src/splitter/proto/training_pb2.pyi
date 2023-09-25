"""
@generated by mypy-protobuf.  Do not edit manually!
isort:skip_file
"""
import builtins
import collections.abc
import google.protobuf.descriptor
import google.protobuf.internal.containers
import google.protobuf.message
import sys

if sys.version_info >= (3, 8):
    import typing as typing_extensions
else:
    import typing_extensions

DESCRIPTOR: google.protobuf.descriptor.FileDescriptor

@typing_extensions.final
class Splitter(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    TRAINING_ID_FIELD_NUMBER: builtins.int
    REPO_URL_FIELD_NUMBER: builtins.int
    training_id: builtins.str
    repo_url: builtins.str
    def __init__(
        self,
        *,
        training_id: builtins.str = ...,
        repo_url: builtins.str = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["repo_url", b"repo_url", "training_id", b"training_id"]) -> None: ...

global___Splitter = Splitter

@typing_extensions.final
class Training(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    TRAINING_ID_FIELD_NUMBER: builtins.int
    MODEL_ID_FIELD_NUMBER: builtins.int
    training_id: builtins.str
    model_id: builtins.str
    def __init__(
        self,
        *,
        training_id: builtins.str = ...,
        model_id: builtins.str = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["model_id", b"model_id", "training_id", b"training_id"]) -> None: ...

global___Training = Training

@typing_extensions.final
class Validation(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    TRAINING_ID_FIELD_NUMBER: builtins.int
    MODELS_SCORES_FIELD_NUMBER: builtins.int
    training_id: builtins.str
    @property
    def models_scores(self) -> google.protobuf.internal.containers.RepeatedCompositeFieldContainer[global___ModelScore]: ...
    def __init__(
        self,
        *,
        training_id: builtins.str = ...,
        models_scores: collections.abc.Iterable[global___ModelScore] | None = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["models_scores", b"models_scores", "training_id", b"training_id"]) -> None: ...

global___Validation = Validation

@typing_extensions.final
class ModelScore(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    MODEL_ID_FIELD_NUMBER: builtins.int
    SCORE_FIELD_NUMBER: builtins.int
    model_id: builtins.str
    score: builtins.int
    def __init__(
        self,
        *,
        model_id: builtins.str = ...,
        score: builtins.int = ...,
    ) -> None: ...
    def ClearField(self, field_name: typing_extensions.Literal["model_id", b"model_id", "score", b"score"]) -> None: ...

global___ModelScore = ModelScore
