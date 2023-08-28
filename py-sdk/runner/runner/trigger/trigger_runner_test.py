from queue import Queue
from unittest.mock import AsyncMock, Mock, call, patch

import pytest
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext

from runner.common.common import Finalizer, Initializer
from runner.trigger.exceptions import UndefinedRunnerFunctionError
from runner.trigger.trigger_runner import ResponseHandler, RunnerFunc, TriggerRunner
from sdk.kai_nats_msg_pb2 import KaiNatsMessage
from sdk.kai_sdk import KaiSDK
from sdk.metadata.metadata import Metadata


@pytest.fixture(scope="function")
async def m_sdk() -> KaiSDK:
    nc = AsyncMock(spec=NatsClient)
    js = Mock(spec=JetStreamContext)
    request_msg = KaiNatsMessage()

    sdk = KaiSDK(nc=nc, js=js)
    sdk.set_request_msg(request_msg)

    return sdk


@pytest.fixture(scope="function")
def m_trigger_runner(m_sdk: KaiSDK) -> TriggerRunner:
    nc = AsyncMock(spec=NatsClient)
    js = Mock(spec=JetStreamContext)

    trigger_runner = TriggerRunner(nc=nc, js=js)

    trigger_runner.response_handler = Mock(spec=ResponseHandler)
    trigger_runner.sdk = m_sdk
    trigger_runner.sdk.metadata = Mock(spec=Metadata)
    trigger_runner.sdk.metadata.get_process = Mock(return_value="test.process")

    return trigger_runner


class MockThread:
    def __init__(self):
        self.start = Mock()
        self.join = Mock()


def test_ok():
    nc = NatsClient()
    js = nc.jetstream()

    runner = TriggerRunner(nc=nc, js=js)

    assert runner.sdk is not None
    assert runner.subscriber is not None


def test_with_initializer_ok(m_trigger_runner):
    m_trigger_runner.with_initializer(AsyncMock(spec=Initializer))

    assert m_trigger_runner.initializer is not None


def test_with_runner_ok(m_trigger_runner):
    m_trigger_runner.with_runner(Mock(spec=ResponseHandler))

    assert m_trigger_runner.runner is not None


def test_with_finalizer_ok(m_trigger_runner):
    m_trigger_runner.with_finalizer(Mock(spec=Finalizer))

    assert m_trigger_runner.finalizer is not None


def test_get_response_channel_ok(m_trigger_runner):
    m_queue = Mock(spec=Queue)
    m_trigger_runner.response_channels = {"test-request-id": m_queue}

    m_trigger_runner.get_response_channel("test-request-id")

    assert m_queue.get.called


@patch("runner.trigger.trigger_runner.Queue", return_value=Mock(spec=Queue))
def test_get_response_channel_not_found_ok(m_queue, m_trigger_runner):
    assert "test-request-id" not in m_trigger_runner.response_channels
    m_trigger_runner.get_response_channel("test-request-id")

    assert "test-request-id" in m_trigger_runner.response_channels
    assert m_queue.called
    assert m_queue.return_value.get.called


@patch("runner.trigger.trigger_runner.get_response_handler", return_value=Mock(spec=ResponseHandler))
@patch("runner.trigger.trigger_runner.Thread", return_value=MockThread())
async def test_run_ok(m_thread, m_response_handler, m_trigger_runner):
    m_trigger_runner.initializer = AsyncMock(spec=Initializer)
    m_trigger_runner.runner = Mock(spec=RunnerFunc)
    m_trigger_runner.finalizer = Mock(spec=Finalizer)

    await m_trigger_runner.run()

    assert m_trigger_runner.initializer.called
    assert m_trigger_runner.initializer.call_args == call(m_trigger_runner.sdk)
    assert m_trigger_runner.response_handler == m_response_handler.return_value
    assert m_thread.call_count == 2
    assert m_thread.call_args_list == [
        call(target=m_trigger_runner.runner, args=(m_trigger_runner, m_trigger_runner.sdk)),
        call(target=m_trigger_runner.subscriber.start, args=()),
    ]
    assert m_thread.return_value.start.call_count == 2
    assert m_thread.return_value.join.call_count == 2
    assert m_trigger_runner.finalizer.called
    assert m_trigger_runner.finalizer.call_args == call(m_trigger_runner.sdk)


async def test_run_undefined_runner_ko(m_trigger_runner):
    with pytest.raises(UndefinedRunnerFunctionError):
        await m_trigger_runner.run()


@patch("runner.trigger.trigger_runner.get_response_handler", return_value=Mock(spec=ResponseHandler))
@patch("runner.trigger.trigger_runner.Thread", return_value=MockThread())
@patch("runner.trigger.trigger_runner.compose_initializer", return_value=AsyncMock(spec=Initializer))
@patch("runner.trigger.trigger_runner.compose_finalizer", return_value=Mock(spec=Finalizer))
async def test_run_undefined_initializer_finalizer_ok(
    m_finalizer, m_initializer, m_thread, m_response_handler, m_trigger_runner
):
    m_trigger_runner.runner = Mock(spec=RunnerFunc)

    await m_trigger_runner.run()

    assert m_trigger_runner.initializer.called
    assert m_trigger_runner.initializer.call_args == call(m_trigger_runner.sdk)
    assert m_trigger_runner.response_handler == m_response_handler.return_value
    assert m_thread.call_count == 2
    assert m_thread.call_args_list == [
        call(target=m_trigger_runner.runner, args=(m_trigger_runner, m_trigger_runner.sdk)),
        call(target=m_trigger_runner.subscriber.start, args=()),
    ]
    assert m_thread.return_value.start.call_count == 2
    assert m_thread.return_value.join.call_count == 2
    assert m_trigger_runner.finalizer.called
    assert m_trigger_runner.finalizer.call_args == call(m_trigger_runner.sdk)
