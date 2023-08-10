from unittest.mock import call, patch

import pytest
from google.protobuf import wrappers_pb2 as wrappers
from google.protobuf.any_pb2 import Any
from google.protobuf.message import Message
from kai_nats_msg_pb2 import KaiNatsMessage, MessageType
from messaging.messaging import Messaging
from mock import AsyncMock, Mock
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v

TEST_CHANNEL = "subscription.test"


@pytest.fixture(scope="function")
def m_messaging() -> Messaging:
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


# async def test_publish_msg(m_messaging):
#     request_id = "test_request_id"
#     payload = wrappers.StringValue(value="test_payload")
#     expected_response_msg = KaiNatsMessage(
#         request_id=request_id,
#         from_node="test_process_id",
#         message_type=MessageType.OK,
#         payload=payload,
#     )
#     m_messaging._new_response_msg = Mock(return_value=expected_response_msg)
#     m_messaging._publish_response = AsyncMock()

#     await m_messaging._publish_msg(msg=payload, msg_type=MessageType.OK, request_id=request_id, chan=TEST_CHANNEL)

#     assert m_messaging._new_response_msg.called
#     assert m_messaging._new_response_msg.call_args == call(Message(), None, MessageType.OK)
#     assert m_messaging._publish_response.called
#     assert m_messaging._publish_response.call_args == call(expected_response_msg, TEST_CHANNEL)


# async def test_publish_any(m_messaging):
#     request_id = "test_request_id"
#     payload = Any()
#     expected_response_msg = KaiNatsMessage(
#         request_id=request_id,
#         payload=payload,
#         from_node="test_process_id",
#         message_type=MessageType.OK,
#     )
#     m_messaging._new_response_msg = Mock(return_value=expected_response_msg)
#     m_messaging._publish_response = AsyncMock()

#     await m_messaging._publish_any(payload=payload, msg_type=MessageType.OK, request_id=request_id, chan=TEST_CHANNEL)

#     assert m_messaging._new_response_msg.called
#     assert m_messaging._new_response_msg.call_args == call(Any(), None, MessageType.OK)
#     assert m_messaging._publish_response.called
#     assert m_messaging._publish_response.call_args == call(expected_response_msg, TEST_CHANNEL)

# @pytest.mark.parametrize(
#         "payload, function",
#         [
#             (Message(), "_publish_msg"),
#             (Any(), "_publish_any"),
#         ],
#     )
# async def test_publish_output(m_messaging, payload, function, expected_call_args):
#     request_id = "test_request_id"
#     expected_response_msg = KaiNatsMessage(
#         request_id=request_id,
#         payload=payload,
#         from_node="test_process_id",
#         message_type=MessageType.OK,
#     )
#     m_messaging._new_response_msg = Mock(return_value=expected_response_msg)
#     m_messaging._publish_response = AsyncMock()

#     await getattr(m_messaging, function)(payload=payload, request_id=request_id, msg_type=MessageType.OK, chan=TEST_CHANNEL)

#     assert m_messaging._new_response_msg.called
#     assert m_messaging._new_response_msg.call_args == call(payload, None, MessageType.OK)
#     assert m_messaging._publish_response.called
#     assert m_messaging._publish_response.call_args == call(expected_response_msg, TEST_CHANNEL)

# async def test_publish_error(m_messaging):
#     m_messaging._publish_response = AsyncMock()
#     v.set("metadata.process_id", "test_process_id")

#     await m_messaging._publish_error(request_id="test_request_id", err_msg="test_error")

#     assert m_messaging._publish_response.called
#     assert m_messaging._publish_response.call_args == call(
#         KaiNatsMessage(
#             request_id="test_request_id",
#             error="test_error",
#             from_mode="test_process_id",
#             message_type=MessageType.ERROR,
#         )
#     )
