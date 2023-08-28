import asyncio
from asyncio import AbstractEventLoop
from threading import Event
from unittest.mock import AsyncMock, Mock, call, patch

import pytest
from google.protobuf.any_pb2 import Any
from nats.aio.client import Client as NatsClient
from nats.aio.client import Msg
from nats.js import JetStreamContext
from nats.js.client import JetStreamContext
from vyper import v

from runner.exit.exceptions import NewRequestMsgError
from runner.exit.exit_runner import ExitRunner
from runner.exit.subscriber import ExitSubscriber
from sdk.kai_nats_msg_pb2 import KaiNatsMessage, MessageType
from sdk.kai_sdk import KaiSDK
from sdk.messaging.messaging_utils import compress
from sdk.metadata.metadata import Metadata

NATS_INPUT = "nats.inputs"


@pytest.fixture(scope="function")
async def m_sdk() -> KaiSDK:
    nc = AsyncMock(spec=NatsClient)
    js = Mock(spec=JetStreamContext)
    request_msg = KaiNatsMessage()

    sdk = KaiSDK(nc=nc, js=js)
    sdk.set_request_msg(request_msg)

    return sdk


@pytest.fixture(scope="function")
def m_exit_runner(m_sdk: KaiSDK) -> ExitRunner:
    nc = AsyncMock(spec=NatsClient)
    js = Mock(spec=JetStreamContext)

    exit_runner = ExitRunner(nc=nc, js=js)
    exit_runner.sdk = m_sdk
    exit_runner.sdk.metadata = Mock(spec=Metadata)
    exit_runner.sdk.metadata.get_process = Mock(return_value="test.process")

    return exit_runner


@pytest.fixture(scope="function")
def m_exit_subscriber(m_exit_runner: ExitRunner) -> ExitSubscriber:
    exit_subscriber = ExitSubscriber(m_exit_runner)

    return exit_subscriber


SUBJECT = "test.subject"


@patch("runner.exit.subscriber.asyncio", return_value=Mock(spec=asyncio))
async def test_start_ok(m_asyncio, m_exit_subscriber):
    v.set(NATS_INPUT, [("1", SUBJECT)])
    m_exit_subscriber.exit_runner.js.subscribe = AsyncMock()
    m_exit_subscriber.subscriber_thread_shutdown_event.set = Mock()
    m_exit_subscriber.subscriber_thread_shutdown_event.wait = Mock()
    m_add_signal_h = m_exit_subscriber.loop.add_signal_handler = Mock()

    await m_exit_subscriber.start()

    assert m_exit_subscriber.exit_runner.js.subscribe.called
    assert m_add_signal_h.call_count == 2
    assert not m_exit_subscriber.subscriber_thread_shutdown_event.set.called
    assert m_exit_subscriber.subscriber_thread_shutdown_event.wait.called


@patch("runner.exit.subscriber.Event", return_value=Mock(spec=Event))
@patch("runner.exit.subscriber.asyncio", return_value=Mock(spec=asyncio))
async def test_start_nats_subscribing_ko(m_asyncio, m_shutdown_event, m_exit_subscriber):
    v.set(NATS_INPUT, [("1", SUBJECT)])
    m_exit_subscriber.exit_runner.js.subscribe = AsyncMock(side_effect=Exception("Subscription error"))
    m_add_signal_h = m_exit_subscriber.loop.add_signal_handler = Mock()

    with pytest.raises(SystemExit):
        await m_exit_subscriber.start()

        assert m_shutdown_event.called
        assert m_exit_subscriber.exit_runner.js.subscribe.called
        assert m_exit_subscriber.loop.stop.called
        assert not m_add_signal_h.called
        assert not m_shutdown_event.return_value.set.called
        assert not m_shutdown_event.return_value.wait.called


async def test_shutdown_handler_coro_ok(m_exit_subscriber):
    v.set(NATS_INPUT, [("1", "test.subject1"), ("2", "test.subject2")])
    m_sub = Mock(spec=JetStreamContext.PushSubscription)
    m_sub.unsubscribe.side_effect = [None, None]
    m_subscriptions = [m_sub, m_sub]
    m_exit_subscriber.exit_runner.js.subscribe = AsyncMock(return_value=m_sub)
    m_exit_subscriber.subscriber_thread_shutdown_event.set = Mock()

    await m_exit_subscriber._shutdown_handler_coro(m_subscriptions)

    assert m_subscriptions[0].unsubscribe.called
    assert m_subscriptions[1].unsubscribe.called
    assert m_exit_subscriber.subscriber_thread_shutdown_event.set.called


async def test_shutdown_handler_coro_ko(m_exit_subscriber):
    v.set(NATS_INPUT, [("1", "test.subject1"), ("2", "test.subject2")])
    m_sub = Mock(spec=JetStreamContext.PushSubscription)
    m_sub.unsubscribe.side_effect = [Exception("Unsubscribe error"), None]
    m_subscriptions = [m_sub, m_sub]
    m_exit_subscriber.exit_runner.js.subscribe = AsyncMock(return_value=m_sub)
    m_exit_subscriber.loop = Mock(spec=AbstractEventLoop)
    m_exit_subscriber.subscriber_thread_shutdown_event.set = Mock()

    with pytest.raises(SystemExit):
        await m_exit_subscriber._shutdown_handler_coro(m_subscriptions)

        assert m_subscriptions[0].unsubscribe.called
        assert m_subscriptions[1].unsubscribe.called
        assert m_exit_subscriber.subscriber_thread_shutdown_event.set.called


async def test_shutdown_handler_ok(m_exit_subscriber):
    v.set(NATS_INPUT, [("1", "test.subject3"), ("2", "test.subject4")])
    m_sub = Mock(spec=JetStreamContext.PushSubscription)
    m_sub.unsubscribe.side_effect = [None, None]
    m_subscriptions = [m_sub, m_sub]
    m_exit_subscriber.exit_runner.js.subscribe = AsyncMock(return_value=m_sub)
    m_exit_subscriber.loop = Mock(spec=AbstractEventLoop)
    m_exit_subscriber._shutdown_handler_coro = AsyncMock()

    m_exit_subscriber._shutdown_handler(m_subscriptions)

    assert m_exit_subscriber._shutdown_handler_coro.called


async def test_process_message_ok(m_exit_subscriber):
    request_id = "test_request_id"
    m_msg = Mock(spec=Msg)
    m_msg.ack = AsyncMock()
    expected_response_msg = KaiNatsMessage(
        request_id=request_id,
        from_node="test_process_id",
        message_type=MessageType.OK,
        payload=Any(),
    )
    m_msg.data = expected_response_msg.SerializeToString()
    m_exit_subscriber._new_request_msg = Mock(return_value=expected_response_msg)
    m_exit_subscriber._process_runner_error = AsyncMock()
    m_handler = Mock()
    m_exit_subscriber.exit_runner.sdk.metadata.get_process = Mock(return_value="test_process_id")
    m_exit_subscriber._get_response_handler = m_handler

    await m_exit_subscriber._process_message(m_msg)

    assert m_exit_subscriber._new_request_msg.called
    assert not m_exit_subscriber._process_runner_error.called
    assert m_exit_subscriber.exit_runner.sdk.request_msg == expected_response_msg
    assert m_exit_subscriber._get_response_handler.called
    assert m_handler.called
    assert m_msg.ack.called


async def test_process_message_not_valid_protobuf_ko(m_exit_subscriber):
    request_id = "test_request_id"
    m_msg = Mock(spec=Msg)
    m_msg.ack = AsyncMock()
    expected_response_msg = KaiNatsMessage(
        request_id=request_id,
        from_node="test_process_id",
        message_type=MessageType.OK,
        payload=Any(),
    )
    m_msg.data = expected_response_msg.SerializeToString()
    m_exit_subscriber._new_request_msg = Mock(side_effect=Exception(NewRequestMsgError("New request message error")))
    m_exit_subscriber._process_runner_error = AsyncMock()

    await m_exit_subscriber._process_message(m_msg)

    assert m_exit_subscriber._new_request_msg.called
    assert m_exit_subscriber._process_runner_error.called
    assert not m_msg.ack.called


async def test_process_message_undefined_handler_ko(m_exit_subscriber):
    request_id = "test_request_id"
    m_msg = Mock(spec=Msg)
    m_msg.ack = AsyncMock()
    expected_response_msg = KaiNatsMessage(
        request_id=request_id,
        from_node="test_process_id",
        message_type=MessageType.OK,
        payload=Any(),
    )
    m_msg.data = expected_response_msg.SerializeToString()
    m_exit_subscriber._new_request_msg = Mock(return_value=expected_response_msg)
    m_exit_subscriber._process_runner_error = AsyncMock()
    m_exit_subscriber.exit_runner.sdk.metadata.get_process = Mock(return_value="test_process_id")
    m_exit_subscriber._get_response_handler = Mock(return_value=None)

    await m_exit_subscriber._process_message(m_msg)

    assert m_exit_subscriber._new_request_msg.called
    m_exit_subscriber.exit_runner.sdk.metadata.get_process.called
    assert m_exit_subscriber._get_response_handler.called
    assert m_exit_subscriber._process_runner_error.called
    assert not m_msg.ack.called


async def test_process_message_preprocessor_ko(m_exit_subscriber):
    request_id = "test_request_id"
    m_msg = Mock(spec=Msg)
    m_msg.ack = AsyncMock()
    expected_response_msg = KaiNatsMessage(
        request_id=request_id,
        from_node="test_process_id",
        message_type=MessageType.OK,
        payload=Any(),
    )
    m_msg.data = expected_response_msg.SerializeToString()
    m_exit_subscriber._new_request_msg = Mock(return_value=expected_response_msg)
    m_exit_subscriber._process_runner_error = AsyncMock()
    m_exit_subscriber.exit_runner.sdk.metadata.get_process = Mock(return_value="test_process_id")
    m_exit_subscriber._get_response_handler = Mock()
    m_exit_subscriber.exit_runner.preprocessor = Mock(side_effect=Exception("Preprocessor error"))

    await m_exit_subscriber._process_message(m_msg)

    assert m_exit_subscriber._new_request_msg.called
    m_exit_subscriber.exit_runner.sdk.metadata.get_process.called
    assert m_exit_subscriber._get_response_handler.called
    assert m_exit_subscriber._process_runner_error.called
    assert not m_msg.ack.called


async def test_process_message_handler_ko(m_exit_subscriber):
    request_id = "test_request_id"
    m_msg = Mock(spec=Msg)
    m_msg.ack = AsyncMock()
    expected_response_msg = KaiNatsMessage(
        request_id=request_id,
        from_node="test_process_id",
        message_type=MessageType.OK,
        payload=Any(),
    )
    m_msg.data = expected_response_msg.SerializeToString()
    m_exit_subscriber._new_request_msg = Mock(return_value=expected_response_msg)
    m_exit_subscriber._process_runner_error = AsyncMock()
    m_exit_subscriber.exit_runner.sdk.metadata.get_process = Mock(return_value="test_process_id")
    m_handler = Mock(side_effect=Exception("Handler error"))
    m_exit_subscriber._get_response_handler = Mock(return_value=m_handler)

    await m_exit_subscriber._process_message(m_msg)

    assert m_exit_subscriber._new_request_msg.called
    m_exit_subscriber.exit_runner.sdk.metadata.get_process.called
    assert m_exit_subscriber._get_response_handler.called
    assert m_exit_subscriber._process_runner_error.called
    assert not m_msg.ack.called


async def test_process_message_postprocessor_ko(m_exit_subscriber):
    request_id = "test_request_id"
    m_msg = Mock(spec=Msg)
    m_msg.ack = AsyncMock()
    expected_response_msg = KaiNatsMessage(
        request_id=request_id,
        from_node="test_process_id",
        message_type=MessageType.OK,
        payload=Any(),
    )
    m_msg.data = expected_response_msg.SerializeToString()
    m_exit_subscriber._new_request_msg = Mock(return_value=expected_response_msg)
    m_exit_subscriber._process_runner_error = AsyncMock()
    m_exit_subscriber.exit_runner.sdk.metadata.get_process = Mock(return_value="test_process_id")
    m_handler = Mock()
    m_exit_subscriber._get_response_handler = m_handler
    m_exit_subscriber.exit_runner.postprocessor = Mock(side_effect=Exception("Postprocessor error"))

    await m_exit_subscriber._process_message(m_msg)

    assert m_exit_subscriber._new_request_msg.called
    m_exit_subscriber.exit_runner.sdk.metadata.get_process.called
    assert m_exit_subscriber._get_response_handler.called
    assert m_handler.called
    assert m_exit_subscriber._process_runner_error.called
    assert not m_msg.ack.called


async def test_process_message_ack_ko_ok(m_exit_subscriber):
    request_id = "test_request_id"
    m_msg = Mock(spec=Msg)
    m_msg.ack = AsyncMock(side_effect=Exception("Ack error"))
    expected_response_msg = KaiNatsMessage(
        request_id=request_id,
        from_node="test_process_id",
        message_type=MessageType.OK,
        payload=Any(),
    )
    m_msg.data = expected_response_msg.SerializeToString()
    m_exit_subscriber._new_request_msg = Mock(return_value=expected_response_msg)
    m_exit_subscriber._process_runner_error = AsyncMock()
    m_exit_subscriber.exit_runner.sdk.metadata.get_process = Mock(return_value="test_process_id")
    m_handler = Mock()
    m_exit_subscriber._get_response_handler = m_handler

    await m_exit_subscriber._process_message(m_msg)

    assert m_exit_subscriber._new_request_msg.called
    m_exit_subscriber.exit_runner.sdk.metadata.get_process.called
    assert m_exit_subscriber._get_response_handler.called
    assert m_handler.called
    assert m_msg.ack.called


async def test_process_runner_error_ok(m_exit_subscriber):
    m_msg = Mock(spec=Msg)
    m_msg.data = b"generic error"
    m_msg.ack = AsyncMock()
    m_exit_subscriber.exit_runner.sdk.messaging.send_error = AsyncMock()

    await m_exit_subscriber._process_runner_error(m_msg, Exception("process runner error"), "test_request_id")

    assert m_msg.ack.called
    assert m_exit_subscriber.exit_runner.sdk.messaging.send_error.called
    assert m_exit_subscriber.exit_runner.sdk.messaging.send_error.call_args == call(
        "process runner error", "test_request_id"
    )


async def test_process_runner_error_ack_ko_ok(m_exit_subscriber):
    m_msg = Mock(spec=Msg)
    m_msg.data = b"generic error"
    m_msg.ack = AsyncMock(side_effect=Exception("Ack error"))
    m_exit_subscriber.exit_runner.sdk.messaging.send_error = AsyncMock()

    await m_exit_subscriber._process_runner_error(m_msg, Exception("process runner error ack"), "test_request_id")

    assert m_exit_subscriber.exit_runner.sdk.messaging.send_error.called
    assert m_exit_subscriber.exit_runner.sdk.messaging.send_error.call_args == call(
        "process runner error ack", "test_request_id"
    )


def test_new_request_msg_ok(m_exit_subscriber):
    request_id = "test_request_id"
    expected_response_msg = KaiNatsMessage(
        request_id=request_id,
        from_node="test_process_id",
        message_type=MessageType.OK,
        payload=Any(),
    )
    data = expected_response_msg.SerializeToString()

    result = m_exit_subscriber._new_request_msg(data)

    assert result == expected_response_msg


def test_new_request_msg_compressed_ok(m_exit_subscriber):
    request_id = "test_request_id"
    expected_response_msg = KaiNatsMessage(
        request_id=request_id,
        from_node="test_process_id",
        message_type=MessageType.OK,
        payload=Any(),
    )
    data = expected_response_msg.SerializeToString()
    data = compress(data)

    result = m_exit_subscriber._new_request_msg(data)

    assert result == expected_response_msg


@patch("runner.exit.subscriber.uncompress", side_effect=Exception("Uncompress error"))
def test_new_request_msg_compressed_ko(_, m_exit_subscriber):
    request_id = "test_request_id"
    expected_response_msg = KaiNatsMessage(
        request_id=request_id,
        from_node="test_process_id",
        message_type=MessageType.OK,
        payload=Any(),
    )
    data = expected_response_msg.SerializeToString()
    data = compress(data)

    with pytest.raises(NewRequestMsgError):
        m_exit_subscriber._new_request_msg(data)


def test_get_response_handler_undefined_default_subject_ok(m_exit_subscriber):
    result = m_exit_subscriber._get_response_handler("wrong_subject")

    assert result is None


def test_get_response_handler_default_subject_ok(m_exit_subscriber):
    m_exit_subscriber.exit_runner.response_handlers = {"default": Mock()}

    result = m_exit_subscriber._get_response_handler("wrong_subject")

    assert result == m_exit_subscriber.exit_runner.response_handlers["default"]


def test_get_response_handler_defined_subject_ok(m_exit_subscriber):
    m_exit_subscriber.exit_runner.response_handlers = {SUBJECT: Mock()}

    result = m_exit_subscriber._get_response_handler(SUBJECT)

    assert result == m_exit_subscriber.exit_runner.response_handlers[SUBJECT]
