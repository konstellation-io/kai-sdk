from unittest.mock import AsyncMock, Mock, call, patch

import pytest
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext

from runner.common.common import Finalizer, Handler, Initializer
from runner.exit.exceptions import UndefinedDefaultHandlerFunctionError
from runner.exit.exit_runner import ExitRunner, Postprocessor, Preprocessor
from runner.exit.subscriber import ExitSubscriber
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
def m_exit_runner(m_sdk: KaiSDK) -> ExitRunner:
    nc = AsyncMock(spec=NatsClient)
    js = Mock(spec=JetStreamContext)

    exit_runner = ExitRunner(nc=nc, js=js)

    exit_runner.sdk = m_sdk
    exit_runner.sdk.metadata = Mock(spec=Metadata)
    exit_runner.sdk.metadata.get_process = Mock(return_value="test.process")
    exit_runner.subscriber = Mock(spec=ExitSubscriber)

    return exit_runner


def test_ok():
    nc = NatsClient()
    js = nc.jetstream()

    runner = ExitRunner(nc=nc, js=js)

    assert runner.sdk is not None
    assert runner.subscriber is not None


def test_with_initializer_ok(m_exit_runner):
    m_exit_runner.with_initializer(AsyncMock(spec=Initializer))

    assert m_exit_runner.initializer is not None


def test_with_prepocessor_ok(m_exit_runner):
    m_exit_runner.with_preprocessor(Mock(spec=Preprocessor))

    assert m_exit_runner.preprocessor is not None


def test_with_handler_ok(m_exit_runner):
    m_exit_runner.with_handler(Mock(spec=Handler))

    assert m_exit_runner.response_handlers["default"] is not None


def test_with_custom_handler_ok(m_exit_runner):
    m_exit_runner.with_custom_handler("test-subject", Mock(spec=Handler))

    assert m_exit_runner.response_handlers["test-subject"] is not None


def test_with_postprocessor_ok(m_exit_runner):
    m_exit_runner.with_postprocessor(Mock(spec=Postprocessor))

    assert m_exit_runner.postprocessor is not None


def test_with_finalizer_ok(m_exit_runner):
    m_exit_runner.with_finalizer(Mock(spec=Finalizer))

    assert m_exit_runner.finalizer is not None


async def test_run_ok(m_exit_runner):
    m_exit_runner.initializer = AsyncMock(spec=Initializer)
    m_exit_runner.finalizer = Mock(spec=Finalizer)
    m_exit_runner.with_handler(Mock(spec=Handler))

    await m_exit_runner.run()

    assert m_exit_runner.initializer.called
    assert m_exit_runner.initializer.call_args == call(m_exit_runner.sdk)
    assert m_exit_runner.subscriber.start.called
    assert m_exit_runner.finalizer.called
    assert m_exit_runner.finalizer.call_args == call(m_exit_runner.sdk)


async def test_run_undefined_runner_ko(m_exit_runner):
    with pytest.raises(UndefinedDefaultHandlerFunctionError):
        await m_exit_runner.run()


@patch("runner.exit.exit_runner.compose_initializer", return_value=AsyncMock(spec=Initializer))
@patch("runner.exit.exit_runner.compose_finalizer", return_value=Mock(spec=Finalizer))
async def test_run_undefined_initializer_finalizer_ok(m_finalizer, m_initializer, m_exit_runner):
    m_exit_runner.with_handler(Mock(spec=Handler))

    await m_exit_runner.run()

    assert m_exit_runner.initializer.called
    assert m_exit_runner.initializer.call_args == call(m_exit_runner.sdk)
    assert m_exit_runner.subscriber.start.called
    assert m_exit_runner.finalizer.called
    assert m_exit_runner.finalizer.call_args == call(m_exit_runner.sdk)