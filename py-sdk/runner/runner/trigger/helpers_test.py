import asyncio
from queue import Queue
from unittest.mock import AsyncMock, Mock, call

import pytest
from google.protobuf.any_pb2 import Any
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v

from runner.trigger.helpers import compose_finalizer, compose_initializer, compose_runner, get_response_handler
from runner.trigger.trigger_runner import ResponseHandler, TriggerRunner
from sdk.centralized_config.centralized_config import CentralizedConfig
from sdk.kai_nats_msg_pb2 import KaiNatsMessage
from sdk.kai_sdk import KaiSDK
from sdk.metadata.metadata import Metadata

CENTRALIZED_CONFIG = "centralized_configuration.process.config"
TEST_REQUEST_ID = "test-request-id"


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


class MockEvent:
    def __init__(self):
        self.is_set = Mock()
        self.set = Mock()
        self.wait = Mock()


async def m_user_initializer_awaitable(sdk):
    assert sdk is not None
    await asyncio.sleep(0.00001)


async def m_user_runner_awaitable(runner, sdk):
    assert sdk is not None
    assert runner is not None
    await asyncio.sleep(0.00001)


def m_user_finalizer(sdk):
    assert sdk is not None


async def test_compose_initializer_with_awaitable_ok(m_sdk):
    v.set(CENTRALIZED_CONFIG, {"key": "value"})
    m_sdk.centralized_config = Mock(spec=CentralizedConfig)
    m_sdk.centralized_config.set_config = AsyncMock()

    await compose_initializer(m_user_initializer_awaitable)(m_sdk)

    assert m_sdk.centralized_config.set_config.called
    assert m_sdk.centralized_config.set_config.call_args == call("key", "value")


async def test_compose_initializer_with_none_ok(m_sdk):
    v.set(CENTRALIZED_CONFIG, {"key": "value"})
    m_sdk.centralized_config = Mock(spec=CentralizedConfig)
    m_sdk.centralized_config.set_config = AsyncMock()

    await compose_initializer()(m_sdk)

    assert m_sdk.centralized_config.set_config.called
    assert m_sdk.centralized_config.set_config.call_args == call("key", "value")


async def test_compose_runner_ok(m_sdk):
    await compose_runner(m_user_runner_awaitable)(m_trigger_runner, m_sdk)


def test_get_response_handler_ok(m_sdk):
    m_queue = Mock(spec=Queue)
    m_sdk.get_request_id = Mock(return_value=TEST_REQUEST_ID)
    handlers = {TEST_REQUEST_ID: m_queue}

    get_response_handler(handlers)(m_sdk, Any())

    assert m_queue.put.called


def test_compose_finalizer_ok(m_sdk):
    compose_finalizer(m_user_finalizer)(m_sdk)


def test_compose_finalizer_with_none_ok(m_sdk):
    compose_finalizer()(m_sdk)
