"""
@generated by mypy-protobuf.  Do not edit manually!
isort:skip_file
"""
import builtins
import google.protobuf.any_pb2
import google.protobuf.descriptor
import google.protobuf.internal.enum_type_wrapper
import google.protobuf.message
import sys
import typing

if sys.version_info >= (3, 10):
    import typing as typing_extensions
else:
    import typing_extensions

DESCRIPTOR: google.protobuf.descriptor.FileDescriptor

class _MessageType:
    ValueType = typing.NewType("ValueType", builtins.int)
    V: typing_extensions.TypeAlias = ValueType

class _MessageTypeEnumTypeWrapper(google.protobuf.internal.enum_type_wrapper._EnumTypeWrapper[_MessageType.ValueType], builtins.type):
    DESCRIPTOR: google.protobuf.descriptor.EnumDescriptor
    UNDEFINED: _MessageType.ValueType  # 0
    OK: _MessageType.ValueType  # 1
    ERROR: _MessageType.ValueType  # 2

class MessageType(_MessageType, metaclass=_MessageTypeEnumTypeWrapper): ...

UNDEFINED: MessageType.ValueType  # 0
OK: MessageType.ValueType  # 1
ERROR: MessageType.ValueType  # 2
global___MessageType = MessageType

@typing_extensions.final
class KaiNatsMessage(google.protobuf.message.Message):
    DESCRIPTOR: google.protobuf.descriptor.Descriptor

    REQUEST_ID_FIELD_NUMBER: builtins.int
    PAYLOAD_FIELD_NUMBER: builtins.int
    ERROR_FIELD_NUMBER: builtins.int
    FROM_NODE_FIELD_NUMBER: builtins.int
    MESSAGE_TYPE_FIELD_NUMBER: builtins.int
    request_id: builtins.str
    @property
    def payload(self) -> google.protobuf.any_pb2.Any: ...
    error: builtins.str
    from_node: builtins.str
    message_type: global___MessageType.ValueType
    def __init__(
        self,
        *,
        request_id: builtins.str = ...,
        payload: google.protobuf.any_pb2.Any | None = ...,
        error: builtins.str = ...,
        from_node: builtins.str = ...,
        message_type: global___MessageType.ValueType = ...,
    ) -> None: ...
    def HasField(self, field_name: typing_extensions.Literal["payload", b"payload"]) -> builtins.bool: ...
    def ClearField(self, field_name: typing_extensions.Literal["error", b"error", "from_node", b"from_node", "message_type", b"message_type", "payload", b"payload", "request_id", b"request_id"]) -> None: ...

global___KaiNatsMessage = KaiNatsMessage
