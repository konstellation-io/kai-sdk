import uuid
from dataclasses import dataclass
from typing import Optional

from google.protobuf.any_pb2 import Any
from google.protobuf.message import Message
from kai_nats_msg_pb2 import KaiNatsMessage, MessageType
from loguru import logger
from loguru._logger import Logger
from messaging.exceptions import FailedGettingMaxMessageSizeError, MessageTooLargeError
from messaging.messaging_utils import MessagingUtils, compress, size_in_mb
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v


@dataclass
class Messaging:
    js: JetStreamContext
    nc: NatsClient
    req_msg: KaiNatsMessage
    messaging_utils: MessagingUtils = None
    logger: Logger = logger.bind(context="[MESSAGING]")

    def __post_init__(self):
        self.messaging_utils = MessagingUtils(js=self.js, nc=self.nc)

    def send_output(self, response: Message, chan: Optional[str] = None):
        self._publish_msg(msg=response, msg_type=MessageType.OK, chan=chan)

    def send_output_with_request_id(self, response: Message, request_id: str, chan: Optional[str] = None):
        self._publish_msg(msg=response, msg_type=MessageType.OK, request_id=request_id, chan=chan)

    def send_any(self, response: Any, chan: Optional[str] = None):
        self._publish_any(payload=response, msg_type=MessageType.OK, chan=chan)

    def send_any_with_request_id(self, response: Any, request_id: str, chan: Optional[str] = None):
        self._publish_any(payload=response, msg_type=MessageType.OK, request_id=request_id, chan=chan)

    # TODO: remove this method
    def send_early_reply(self, response: Message, chan: Optional[str] = None):
        self._publish_msg(msg=response, msg_type=MessageType.EARLY_REPLY, chan=chan)

    # TODO: remove this method
    def send_early_exit(self, response: Message, chan: Optional[str] = None):
        self._publish_msg(msg=response, msg_type=MessageType.EARLY_EXIT, chan=chan)

    def get_error_message(self) -> str:
        return self.req_msg.error if self.is_message_error() else ""

    def is_message_ok(self) -> bool:
        return self.req_msg.message_type == MessageType.OK

    def is_message_error(self) -> bool:
        return self.req_msg.message_type == MessageType.ERROR

    def is_message_early_reply(self) -> bool:
        return self.req_msg.message_type == MessageType.EARLY_REPLY

    def is_message_early_exit(self) -> bool:
        return self.req_msg.message_type == MessageType.EARLY_EXIT

    def _publish_msg(self, msg: Message, msg_type: int, request_id: Optional[str] = None, chan: Optional[str] = None):
        try:
            payload = Any()
            payload.Pack(msg)
        except Exception as e:
            self.logger.debug(f"failed packing message: {e}")
            return

        if not request_id:
            request_id = str(uuid.uuid4())

        response_msg = self._new_response_msg(payload, request_id, msg_type)
        self._publish_response(response_msg, chan)

    def _publish_any(self, payload: Any, msg_type: int, request_id: Optional[str] = None, chan: Optional[str] = None):
        if not request_id:
            request_id = str(uuid.uuid4())

        response_msg = self._new_response_msg(payload, request_id, msg_type)
        self._publish_response(response_msg, chan)

    def _publish_error(self, request_id: str, err_msg: str):
        response_msg = KaiNatsMessage(
            request_id=request_id,
            error=err_msg,
            from_node=v.get("metadata.process_id"),
            message_type=MessageType.ERROR,
        )
        self._publish_response(response_msg)

    def _new_response_msg(self, payload: Any, request_id: str, msg_type: int) -> KaiNatsMessage:
        self.logger.info(f"preparing response message of type {msg_type} and request_id {request_id}...")
        return KaiNatsMessage(
            RequestId=request_id,
            Payload=payload,
            FromNode=v.get("metadata.process_id"),
            MessageType=msg_type,
        )

    async def _publish_response(self, response_msg: KaiNatsMessage, chan: Optional[str] = None):
        output_subject = self._get_output_subject(chan)

        try:
            output_msg = response_msg.SerializeToString()
        except Exception as e:
            self.logger.debug(f"failed serializing response message: {e}")
            return

        try:
            output_msg = self._prepare_output_message(output_msg)
        except (FailedGettingMaxMessageSizeError, MessageTooLargeError) as e:
            self.logger.debug(f"failed preparing output message: {e}")
            return

        self.logger.info(f"publishing response to subject {output_subject}...")

        try:
            await self.js.publish(output_subject, output_msg)
        except Exception as e:
            self.logger.debug(f"failed publishing response: {e}")
            return

    def _get_output_subject(self, chan: Optional[str] = None) -> str:
        output_subject = v.get("nats.output")
        return f"{output_subject}.{chan}" if chan else output_subject

    def _prepare_output_message(self, msg: bytes) -> bytes | Exception:
        max_size = self.message_utils.get_max_message_size()

        if len(msg) <= max_size:
            return msg

        self.logger.info("message exceeds maximum size allowed! compressing data...")
        out_msg = compress(msg)

        len_out_msg = len(out_msg)
        if len_out_msg > max_size:
            self.logger.warning(f"compressed message size: {len_out_msg} exceeds maximum allowed size: {max_size}")
            raise MessageTooLargeError(size_in_mb(len_out_msg), size_in_mb(max_size))

        self.logger.info(
            f"message compressed! original message size: {len(msg)} - compressed message size: {len_out_msg}"
        )

        return out_msg
