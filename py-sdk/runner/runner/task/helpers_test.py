import asyncio
from typing import Callable
from unittest.mock import AsyncMock, Mock, call

import pytest
from google.protobuf.any_pb2 import Any
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v

from runner.task.helpers import (
    compose_finalizer,
    compose_handler,
    compose_initializer,
    compose_postprocessor,
    compose_preprocessor,
)
from sdk.centralized_config.centralized_config import CentralizedConfig
from sdk.kai_nats_msg_pb2 import KaiNatsMessage
from sdk.kai_sdk import KaiSDK

CENTRALIZED_CONFIG = "centralized_configuration.process.config"


@pytest.fixture(scope="function")
async def m_sdk() -> KaiSDK:
    nc = AsyncMock(spec=NatsClient)
    js = Mock(spec=JetStreamContext)
    request_msg = KaiNatsMessage()

    sdk = KaiSDK(nc=nc, js=js)
    sdk.set_request_msg(request_msg)

    return sdk


async def m_user_initializer_awaitable(sdk):
    assert sdk is not None
    await asyncio.sleep(0.00001)


async def m_user_preprocessor_awaitable(sdk, response):
    assert sdk is not None
    assert response is not None
    await asyncio.sleep(0.00001)


def m_user_handler(sdk, response):
    assert sdk is not None
    assert response is not None


async def m_user_postprocessor_awaitable(sdk, response):
    assert sdk is not None
    assert response is not None
    await asyncio.sleep(0.00001)


def m_user_finalizer(sdk):
    assert sdk is not None


async def test_compose_initializer_with_awaitable_ok(m_sdk):
    v.set(CENTRALIZED_CONFIG, {"key": "value"})
    m_sdk.centralized_config = Mock(spec=CentralizedConfig)
    m_sdk.centralized_config.set_config = AsyncMock()
    initializer: Callable = compose_initializer(m_user_initializer_awaitable)

    await initializer(m_sdk)

    assert m_sdk.centralized_config.set_config.called
    assert m_sdk.centralized_config.set_config.call_args == call("key", "value")


async def test_compose_initializer_with_none_ok(m_sdk):
    v.set(CENTRALIZED_CONFIG, {"key": "value"})
    m_sdk.centralized_config = Mock(spec=CentralizedConfig)
    m_sdk.centralized_config.set_config = AsyncMock()
    initializer: Callable = compose_initializer()

    await initializer(m_sdk)

    assert m_sdk.centralized_config.set_config.called
    assert m_sdk.centralized_config.set_config.call_args == call("key", "value")


async def test_compose_preprocessor_ok(m_sdk):
    preprocessor: Callable = compose_preprocessor(m_user_preprocessor_awaitable)
    await preprocessor(m_sdk, Any())


async def test_compose_handler_ok(m_sdk):
    handler: Callable = compose_handler(m_user_handler)
    await handler(m_sdk, Any())


async def test_compose_postprocessor_ok(m_sdk):
    postprocessor: Callable = compose_postprocessor(m_user_postprocessor_awaitable)
    await postprocessor(m_sdk, Any())


async def test_compose_finalizer_ok(m_sdk):
    finalizer: Callable = compose_finalizer(m_user_finalizer)
    await finalizer(m_sdk)


async def test_compose_finalizer_with_none_ok(m_sdk):
    finalizer: Callable = compose_finalizer()
    await finalizer(m_sdk)
