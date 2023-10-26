import asyncio
from unittest.mock import AsyncMock, Mock, patch

import pytest
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext

from runner.common.common import Finalizer, Initializer
from runner.trigger.exceptions import UndefinedRunnerFunctionError
from runner.trigger.trigger_runner import ResponseHandler, TriggerRunner
from sdk.kai_nats_msg_pb2 import KaiNatsMessage
from sdk.kai_sdk import KaiSDK
from sdk.metadata.metadata import Metadata
from sdk.persistent_storage.persistent_storage import PersistentStorage


@pytest.fixture(scope="function")
@patch.object(PersistentStorage, "__new__", return_value=Mock(spec=PersistentStorage))
async def m_sdk(_) -> KaiSDK:
    nc = AsyncMock(spec=NatsClient)
    js = Mock(spec=JetStreamContext)
    request_msg = KaiNatsMessage()

    sdk = KaiSDK(nc=nc, js=js)
    sdk.set_request_msg(request_msg)

    return sdk


@pytest.fixture(scope="function")
@patch.object(PersistentStorage, "__new__", return_value=Mock(spec=PersistentStorage))
def m_trigger_runner(_, m_sdk: KaiSDK) -> TriggerRunner:
    nc = AsyncMock(spec=NatsClient)
    js = Mock(spec=JetStreamContext)

    trigger_runner = TriggerRunner(nc=nc, js=js)

    trigger_runner.response_handler = Mock(spec=ResponseHandler)
    trigger_runner.sdk = m_sdk
    trigger_runner.sdk.metadata = Mock(spec=Metadata)
    trigger_runner.sdk.metadata.get_process = Mock(return_value="test.process")

    return trigger_runner


class MockAsyncio:
    def __init__(self) -> None:
        self.get_event_loop = Mock()
        self.run_in_executor = Mock()
        self.gather = AsyncMock()


@patch.object(PersistentStorage, "__new__", return_value=Mock(spec=PersistentStorage))
def test_ok(_):
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


async def test_get_response_channel_ok(m_trigger_runner):
    m_queue = AsyncMock(spec=asyncio.Queue)
    m_trigger_runner.response_channels = {"test-request-id": m_queue}

    await m_trigger_runner.get_response_channel("test-request-id")

    assert m_queue.get.called


@patch("runner.trigger.trigger_runner.Queue", return_value=AsyncMock(spec=asyncio.Queue))
async def test_get_response_channel_not_found_ok(m_queue, m_trigger_runner):
    assert "test-request-id" not in m_trigger_runner.response_channels
    await m_trigger_runner.get_response_channel("test-request-id")

    assert "test-request-id" in m_trigger_runner.response_channels
    assert m_queue.called
    assert m_queue.return_value.get.called


async def test_run_undefined_runner_function_ko(m_trigger_runner):
    with pytest.raises(UndefinedRunnerFunctionError):
        await m_trigger_runner.run()
