import asyncio
from unittest.mock import AsyncMock, Mock, call

from vyper import v

from runner.task.helpers import (
    compose_finalizer,
    compose_handler,
    compose_initializer,
    compose_postprocessor,
    compose_preprocessor,
)
from sdk.centralized_config.centralized_config import CentralizedConfig
from sdk.kai_sdk import KaiSDK

CENTRALIZED_CONFIG = "centralized_configuration.process.config"


async def mock_user_initializer_awaitable(sdk):
    assert sdk is not None
    await asyncio.sleep(0.00001)


def mock_user_preprocessor(sdk, response):
    assert sdk is not None
    assert response is not None


def mock_user_handler(sdk, response):
    assert sdk is not None
    assert response is not None


def mock_user_postprocessor(sdk, response):
    assert sdk is not None
    assert response is not None


def mock_user_finalizer(sdk):
    assert sdk is not None


async def test_compose_initializer_with_awaitable_ok():
    v.set(CENTRALIZED_CONFIG, {"key": "value"})
    sdk = Mock(spec=KaiSDK)
    sdk.centralized_config = Mock(spec=CentralizedConfig)
    sdk.centralized_config.set_config = AsyncMock()
    initializer_func = compose_initializer(mock_user_initializer_awaitable)

    await initializer_func(sdk)

    assert sdk.centralized_config.set_config.called
    assert sdk.centralized_config.set_config.call_args == call("key", "value")


async def test_compose_initializer_with_none_ok():
    v.set(CENTRALIZED_CONFIG, {"key": "value"})
    sdk = Mock(spec=KaiSDK)
    sdk.centralized_config = Mock(spec=CentralizedConfig)
    sdk.centralized_config.set_config = AsyncMock()
    initializer_func = compose_initializer()

    await initializer_func(sdk)

    assert sdk.centralized_config.set_config.called
    assert sdk.centralized_config.set_config.call_args == call("key", "value")


async def test_compose_preprocessor_ok():
    sdk = Mock(spec=KaiSDK)
    preprocessor_func = compose_preprocessor(mock_user_preprocessor)

    preprocessor_func(sdk, "response")


async def test_compose_handler_ok():
    sdk = Mock(spec=KaiSDK)
    handler_func = compose_handler(mock_user_handler)

    handler_func(sdk, "response")


async def test_compose_postprocessor_ok():
    sdk = Mock(spec=KaiSDK)
    postprocessor_func = compose_postprocessor(mock_user_postprocessor)

    postprocessor_func(sdk, "response")


async def test_compose_finalizer_ok():
    sdk = Mock(spec=KaiSDK)
    finalizer_func = compose_finalizer(mock_user_finalizer)

    finalizer_func(sdk)


async def test_compose_finalizer_with_none_ok():
    sdk = Mock(spec=KaiSDK)
    finalizer_func = compose_finalizer()

    finalizer_func(sdk)
