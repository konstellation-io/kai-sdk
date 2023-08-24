import asyncio
from signal import signal
from threading import Event
from unittest.mock import AsyncMock, Mock, patch

import pytest
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v

from runner.exit.exit_runner import ExitRunner
from runner.exit.subscriber import ExitSubscriber
from sdk.kai_nats_msg_pb2 import KaiNatsMessage
from sdk.kai_sdk import KaiSDK
from sdk.metadata.metadata import Metadata

NATS_INPUT = "nats.inputs"


@pytest.fixture(scope="function")
async def m_sdk() -> KaiSDK:
    nc = AsyncMock(spec=NatsClient)
    js = Mock(spec=JetStreamContext)
    req_msg = KaiNatsMessage()

    sdk = KaiSDK(nc=nc, js=js)
    sdk.set_request_message(req_msg)

    return sdk


@pytest.fixture(scope="function")
def m_exit_runner(m_sdk) -> ExitRunner:
    nc = AsyncMock(spec=NatsClient)
    js = Mock(spec=JetStreamContext)

    exit_runner = ExitRunner(nc=nc, js=js)
    exit_runner.sdk = m_sdk
    exit_runner.sdk.metadata = Mock(spec=Metadata)
    exit_runner.sdk.metadata.get_process = Mock(return_value="test.process")

    return exit_runner


@patch("runner.exit.subscriber.Event", return_value=Mock(spec=Event))
@patch("runner.exit.subscriber.asyncio", return_value=Mock(spec=asyncio))
async def test_start_ok(m_asyncio, m_shutdown_event, m_exit_runner):
    v.set(NATS_INPUT, [("1", "test.subject")])
    m_exit_runner.js.subscribe = AsyncMock()
    m_add_signal_h = m_asyncio.get_event_loop.return_value.add_signal_handler = Mock()

    s = ExitSubscriber(m_exit_runner)
    await s.start()

    assert m_shutdown_event.called
    assert m_exit_runner.js.subscribe.called
    assert m_add_signal_h.call_count == 2
    assert not m_shutdown_event.return_value.set.called
    assert m_shutdown_event.return_value.wait.called


@patch("runner.exit.subscriber.Event", return_value=Mock(spec=Event))
@patch("runner.exit.subscriber.asyncio", return_value=Mock(spec=asyncio))
async def test_start_nats_subscribing_ko(m_asyncio, m_shutdown_event, m_exit_runner):
    v.set(NATS_INPUT, [("1", "test.subject")])
    m_exit_runner.js.subscribe = AsyncMock(side_effect=Exception("Subscription error"))
    m_add_signal_h = m_asyncio.get_event_loop.return_value.add_signal_handler = Mock()

    with pytest.raises(SystemExit):
        s = ExitSubscriber(m_exit_runner)
        await s.start()

        assert m_shutdown_event.called
        assert m_exit_runner.js.subscribe.called
        assert m_asyncio.return_value.get_event_loop.return_value.stop.called
        assert not m_add_signal_h.called
        assert not m_shutdown_event.return_value.set.called
        assert not m_shutdown_event.return_value.wait.called


#     def side_effect_shutdown_handler(sig, frame):
#         instance.subscriber_thread_shutdown_event.set()

#     m_create_task.side_effect = side_effect_shutdown_handler
