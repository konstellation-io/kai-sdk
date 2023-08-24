import asyncio
from signal import signal
from threading import Event
from unittest.mock import AsyncMock, Mock, patch

import pytest
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v

from runner.task.subscriber import TaskSubscriber
from runner.task.task_runner import TaskRunner
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
def m_task_runner(m_sdk) -> TaskRunner:
    nc = AsyncMock(spec=NatsClient)
    js = Mock(spec=JetStreamContext)

    task_runner = TaskRunner(nc=nc, js=js)
    task_runner.sdk = m_sdk
    task_runner.sdk.metadata = Mock(spec=Metadata)
    task_runner.sdk.metadata.get_process = Mock(return_value="test.process")

    return task_runner


@patch("runner.task.subscriber.Event", return_value=Mock(spec=Event))
@patch("runner.task.subscriber.asyncio", return_value=Mock(spec=asyncio))
@patch("runner.task.subscriber.signal", return_value=Mock(spec=signal))
async def test_start_ok(m_signal_task, _, m_shutdown_event, m_task_runner):
    v.set(NATS_INPUT, [("1", "test.subject")])
    m_task_runner.js.subscribe = AsyncMock()

    s = TaskSubscriber(m_task_runner)
    await s.start()

    assert m_shutdown_event.called
    assert m_task_runner.js.subscribe.called
    assert m_signal_task.call_count == 2
    assert not m_shutdown_event.return_value.set.called
    assert m_shutdown_event.return_value.wait.called


@patch("runner.task.subscriber.Event", return_value=Mock(spec=Event))
@patch("runner.task.subscriber.asyncio", return_value=Mock(spec=asyncio))
@patch("runner.task.subscriber.asyncio.get_event_loop.add_signal_handler")
async def test_start_nats_subscribing_ko(m_signal_task, m_asyncio, m_shutdown_event, m_task_runner):
    v.set(NATS_INPUT, [("1", "test.subject")])
    m_task_runner.js.subscribe = AsyncMock(side_effect=Exception("Subscription error"))

    with pytest.raises(SystemExit):
        s = TaskSubscriber(m_task_runner)
        await s.start()

        assert m_shutdown_event.called
        assert m_task_runner.js.subscribe.called
        assert m_asyncio.return_value.get_event_loop.return_value.stop.called
        assert m_signal_task.call_count == 0
        assert not m_shutdown_event.return_value.set.called
        assert not m_shutdown_event.return_value.wait.called


#     def side_effect_shutdown_handler(sig, frame):
#         instance.subscriber_thread_shutdown_event.set()

#     m_create_task.side_effect = side_effect_shutdown_handler
