# Generated by the protocol buffer compiler.  DO NOT EDIT!
# sources: kai_nats_msg.proto
# plugin: python-betterproto
from dataclasses import dataclass

import betterproto
from google.protobuf import any_pb2 as protobuf


class MessageType(betterproto.Enum):
    UNDEFINED = 0
    OK = 1
    ERROR = 2
    EARLY_REPLY = 3
    EARLY_EXIT = 4


@dataclass
class KaiNatsMessage(betterproto.Message):
    request_id: str = betterproto.string_field(1)
    payload: protobuf.Any = betterproto.message_field(2)
    error: str = betterproto.string_field(3)
    from_node: str = betterproto.string_field(4)
    message_type: "MessageType" = betterproto.enum_field(5)
