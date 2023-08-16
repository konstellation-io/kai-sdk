from unittest.mock import call, patch

import pytest
from google.protobuf import wrappers_pb2 as wrappers
from google.protobuf.any_pb2 import Any
from google.protobuf.message import Message
from mock import AsyncMock, Mock
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v

from sdk.kai_nats_msg_pb2 import KaiNatsMessage, MessageType
from sdk.messaging.exceptions import FailedGettingMaxMessageSizeError, MessageTooLargeError
from sdk.messaging.messaging import Messaging
from sdk.messaging.messaging_utils import compress, is_compressed

NATS_OUTPUT = "subscription.output"
TEST_CHANNEL = "subscription.test"
ANY_BYTE = b"any"


@pytest.fixture(scope="function")
def m_messaging() -> Messaging:
    v.set("nats.output", NATS_OUTPUT)
    v.set("metadata.process_id", "test_process_id")
    nc = AsyncMock(spec=NatsClient)
    js = AsyncMock(spec=JetStreamContext)
    req_msg = Mock(spec=KaiNatsMessage)

    messaging = Messaging(nc=nc, js=js, req_msg=req_msg)

    return messaging


def test_ok():
    nc = NatsClient()
    js = nc.jetstream()
    req_msg = KaiNatsMessage()

    messaging = Messaging(nc=nc, js=js, req_msg=req_msg)

    assert messaging is not None
    assert messaging.js is not None
    assert messaging.nc is not None
    assert messaging.req_msg is not None
    assert messaging.logger is not None
    assert messaging.messaging_utils is not None


async def test_send_output(m_messaging):
    m_messaging._publish_msg = AsyncMock()
    response = Message()

    await m_messaging.send_output(response=response, chan=TEST_CHANNEL)

    assert m_messaging._publish_msg.called
    assert m_messaging._publish_msg.call_args == call(msg=response, msg_type=MessageType.OK, chan=TEST_CHANNEL)


async def test_send_output_with_request_id(m_messaging):
    m_messaging._publish_msg = AsyncMock()
    response = Message()
    request_id = "test_request_id"

    await m_messaging.send_output_with_request_id(response=response, request_id=request_id, chan=TEST_CHANNEL)

    assert m_messaging._publish_msg.called
    assert m_messaging._publish_msg.call_args == call(
        msg=response, msg_type=MessageType.OK, request_id=request_id, chan=TEST_CHANNEL
    )


async def test_send_any(m_messaging):
    m_messaging._publish_any = AsyncMock()
    response = Any()

    await m_messaging.send_any(response=response, chan=TEST_CHANNEL)

    assert m_messaging._publish_any.called
    assert m_messaging._publish_any.call_args == call(payload=response, msg_type=MessageType.OK, chan=TEST_CHANNEL)


async def test_send_any_with_request_id(m_messaging):
    m_messaging._publish_any = AsyncMock()
    response = Any()
    request_id = "test_request_id"

    await m_messaging.send_any_with_request_id(response=response, request_id=request_id, chan=TEST_CHANNEL)

    assert m_messaging._publish_any.called
    assert m_messaging._publish_any.call_args == call(
        payload=response, msg_type=MessageType.OK, request_id=request_id, chan=TEST_CHANNEL
    )


async def test_send_early_reply(m_messaging):
    m_messaging._publish_msg = AsyncMock()
    response = Message()

    await m_messaging.send_early_reply(response=response, chan=TEST_CHANNEL)

    assert m_messaging._publish_msg.called
    assert m_messaging._publish_msg.call_args == call(msg=response, msg_type=MessageType.EARLY_REPLY, chan=TEST_CHANNEL)


async def test_send_early_exit(m_messaging):
    m_messaging._publish_msg = AsyncMock()
    response = Message()

    await m_messaging.send_early_exit(response=response, chan=TEST_CHANNEL)

    assert m_messaging._publish_msg.called
    assert m_messaging._publish_msg.call_args == call(msg=response, msg_type=MessageType.EARLY_EXIT, chan=TEST_CHANNEL)


def test_get_error_message(m_messaging):
    m_messaging.req_msg.message_type = MessageType.ERROR
    m_messaging.req_msg.error = "test_error"

    message = m_messaging.get_error_message()

    assert message == "test_error"


@pytest.mark.parametrize(
    "message_type, function, expected_result",
    [
        (MessageType.OK, "is_message_ok", True),
        (MessageType.ERROR, "is_message_error", True),
        (MessageType.EARLY_REPLY, "is_message_early_reply", True),
        (MessageType.EARLY_EXIT, "is_message_early_exit", True),
    ],
)
def test_is_message_ok(m_messaging, message_type, function, expected_result):
    m_messaging.req_msg.message_type = message_type

    is_message = getattr(m_messaging, function)()

    assert is_message == expected_result


@patch.object(Any, "Pack", return_value=Any())
async def test__publish_msg_ok(_, m_messaging):
    request_id = "test_request_id"
    msg = Mock(spec=Message)
    expected_response_msg = KaiNatsMessage(
        request_id=request_id,
        from_node="test_process_id",
        message_type=MessageType.OK,
        payload=Any(),
    )
    m_messaging._publish_response = AsyncMock()

    await m_messaging._publish_msg(msg=msg, msg_type=MessageType.OK, request_id=request_id, chan=TEST_CHANNEL)

    assert m_messaging._publish_response.called
    assert m_messaging._publish_response.call_args == call(expected_response_msg, TEST_CHANNEL)


@patch.object(Any, "Pack", side_effect=Exception)
async def test__publish_msg_packing_message_ko(_, m_messaging):
    message = Any()
    m_messaging._new_response_msg = Mock()
    m_messaging._publish_response = AsyncMock()

    await m_messaging._publish_msg(msg=message, msg_type=MessageType.OK)

    assert not m_messaging._new_response_msg.called
    assert not m_messaging._publish_response.called


async def test__publish_any_ok(m_messaging):
    request_id = "test_request_id"
    payload = Any()
    expected_response_msg = KaiNatsMessage(
        request_id=request_id,
        payload=payload,
        from_node="test_process_id",
        message_type=MessageType.OK,
    )
    m_messaging._publish_response = AsyncMock()

    await m_messaging._publish_any(payload=payload, msg_type=MessageType.OK, request_id=request_id, chan=TEST_CHANNEL)

    assert m_messaging._publish_response.called
    assert m_messaging._publish_response.call_args == call(expected_response_msg, TEST_CHANNEL)


async def test__publish_error_ok(m_messaging):
    m_messaging._publish_response = AsyncMock()

    await m_messaging._publish_error(request_id="test_request_id", err_msg="test_error")

    assert m_messaging._publish_response.called
    assert m_messaging._publish_response.call_args == call(
        KaiNatsMessage(
            request_id="test_request_id",
            error="test_error",
            from_node="test_process_id",
            message_type=MessageType.ERROR,
        )
    )


def test__new_response_msg_ok(m_messaging):
    request_id = "test_request_id"
    payload = Any()

    response = m_messaging._new_response_msg(payload, request_id, msg_type=MessageType.OK)

    assert response == KaiNatsMessage(
        request_id=request_id,
        payload=payload,
        from_node="test_process_id",
        message_type=MessageType.OK,
    )


async def test__publish_response_ok(m_messaging):
    message = KaiNatsMessage()
    bytes_message = message.SerializeToString()
    m_messaging._prepare_output_message = AsyncMock(return_value=bytes_message)

    await m_messaging._publish_response(message)

    assert m_messaging.js.publish.called
    assert m_messaging.js.publish.call_args == call(NATS_OUTPUT, bytes_message)


async def test__publish_response_with_channel_ok(m_messaging):
    message = KaiNatsMessage()
    bytes_message = message.SerializeToString()
    m_messaging._prepare_output_message = AsyncMock(return_value=bytes_message)

    await m_messaging._publish_response(message, chan=TEST_CHANNEL)

    assert m_messaging.js.publish.called
    assert m_messaging.js.publish.call_args == call(f"{NATS_OUTPUT}.{TEST_CHANNEL}", bytes_message)


@patch.object(KaiNatsMessage, "SerializeToString", side_effect=Exception)
async def test__publish_response_serializing_message_ko(_, m_messaging):
    message = KaiNatsMessage()

    await m_messaging._publish_response(message)

    assert not m_messaging.js.publish.called


async def test__publish_response_preparing_output_message_ko(m_messaging):
    message = KaiNatsMessage()
    m_messaging._prepare_output_message = AsyncMock(side_effect=MessageTooLargeError(0, 0))

    await m_messaging._publish_response(message)

    assert not m_messaging.js.publish.called


async def test__publish_response_publishing_message_ko(m_messaging):
    message = KaiNatsMessage()
    m_messaging._prepare_output_message = AsyncMock(return_value=b"any")
    m_messaging.js.publish = AsyncMock(side_effect=Exception)

    await m_messaging._publish_response(message)

    assert m_messaging.js.publish.called


def test__get_output_subject_ok(m_messaging):
    channel = "test_channel"

    output_subject = m_messaging._get_output_subject()
    output_subject_with_channel = m_messaging._get_output_subject(channel)

    assert output_subject == NATS_OUTPUT
    assert output_subject_with_channel == f"{NATS_OUTPUT}.{channel}"


async def test__prepare_output_message_ok(m_messaging):
    m_messaging.messaging_utils.get_max_message_size = AsyncMock(return_value=100)

    result = await m_messaging._prepare_output_message(ANY_BYTE)

    assert result == ANY_BYTE


async def test__prepare_output_message_compressed_ok(m_messaging):
    m_messaging.messaging_utils.get_max_message_size = AsyncMock(return_value=30)

    result = await m_messaging._prepare_output_message(b"any" * 30)

    assert is_compressed(result)
    assert result == compress(b"any" * 30)


async def test__prepare_output_message_getting_max_message_size_ko(m_messaging):
    m_messaging.messaging_utils.get_max_message_size = AsyncMock(side_effect=FailedGettingMaxMessageSizeError)

    with pytest.raises(FailedGettingMaxMessageSizeError):
        await m_messaging._prepare_output_message(ANY_BYTE)


async def test__prepare_output_message_too_large_ko(m_messaging):
    m_messaging.messaging_utils.get_max_message_size = AsyncMock(return_value=0)

    with pytest.raises(MessageTooLargeError):
        await m_messaging._prepare_output_message(ANY_BYTE)