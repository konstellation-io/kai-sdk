import pytest
from mock import AsyncMock, Mock, patch
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from nats.js.kv import KeyValue
from vyper import v

from runner.exit.exit_runner import ExitRunner
from runner.runner import Runner
from runner.task.task_runner import TaskRunner
from runner.trigger.trigger_runner import TriggerRunner
from sdk.centralized_config.centralized_config import CentralizedConfig
from sdk.kai_nats_msg_pb2 import KaiNatsMessage
from sdk.kai_sdk import KaiSDK
from sdk.messaging.messaging import Messaging
from sdk.metadata.metadata import Metadata
from sdk.object_store.object_store import ObjectStore
from sdk.path_utils.path_utils import PathUtils

PRODUCT_BUCKET = "centralized_configuration.product.bucket"
WORKFLOW_BUCKET = "centralized_configuration.workflow.bucket"
PROCESS_BUCKET = "centralized_configuration.process.bucket"
NATS_OBJECT_STORE = "nats.object_store"
NATS_URL = "nats.url"


@pytest.fixture(scope="function")
def m_runner() -> Runner:
    nc = AsyncMock(spec=NatsClient)
    js = Mock(spec=JetStreamContext)
    v.set(NATS_URL, "test_url")
    v.set("APP_CONFIG_PATH", "test_path")
    v.set("runner.logger.output_paths", ["stdout"])
    v.set("runner.logger.error_output_paths", ["stderr"])
    v.set("runner.logger.level", "INFO")

    runner = Runner(nc=nc)
    runner.js = js

    return runner


@patch.object(
    CentralizedConfig,
    "_init_kv_stores",
    return_value=(AsyncMock(spec=KeyValue), AsyncMock(spec=KeyValue), AsyncMock(spec=KeyValue)),
)
async def test_sdk_import_ok(_):
    nc = NatsClient()
    js = nc.jetstream()
    req_msg = KaiNatsMessage()
    v.set(NATS_OBJECT_STORE, None)
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    sdk = KaiSDK(nc=nc, js=js)
    await sdk.initialize()
    sdk.set_request_message(req_msg)

    assert isinstance(sdk.metadata, Metadata)
    assert isinstance(sdk.messaging, Messaging)
    assert isinstance(sdk.object_store, ObjectStore)
    assert isinstance(sdk.centralized_config, CentralizedConfig)
    assert isinstance(sdk.path_utils, PathUtils)
    assert sdk.nc is not None
    assert sdk.js is not None
    assert sdk.req_msg == req_msg
    assert sdk.logger is not None
    assert sdk.metadata is not None
    assert sdk.messaging is not None
    assert sdk.messaging.req_msg == req_msg
    assert sdk.object_store.object_store_name == ""
    assert sdk.object_store.object_store is None
    assert sdk.centralized_config is not None
    assert isinstance(sdk.centralized_config.process_kv, KeyValue)
    assert isinstance(sdk.centralized_config.product_kv, KeyValue)
    assert isinstance(sdk.centralized_config.workflow_kv, KeyValue)
    assert sdk.path_utils is not None
    assert getattr(sdk, "measurements", None) is None
    assert getattr(sdk, "storage", None) is None


async def test_runner_ok():
    nc = NatsClient()
    v.set(NATS_URL, "test_url")
    v.set("APP_CONFIG_PATH", "test_path")
    v.set("runner.logger.output_paths", ["stdout"])
    v.set("runner.logger.error_output_paths", ["stderr"])
    v.set("runner.logger.level", "INFO")

    runner = Runner(nc=nc)

    assert runner.nc is not None
    assert getattr(runner, "js", None) is None
    assert runner.logger is not None


async def test_runner_initialize_ok():
    nc = AsyncMock(spec=NatsClient)
    nc.connect.return_value = None
    v.set(NATS_URL, "test_url")
    m_js = Mock(spec=JetStreamContext)
    nc.jetstream = Mock(return_value=m_js)

    runner = Runner(nc=nc)
    await runner.initialize()

    assert runner.nc is nc
    assert runner.js is m_js


async def test_runner_initialize_nats_ko():
    nc = AsyncMock(spec=NatsClient)
    nc.connect.side_effect = Exception("test exception")
    v.set(NATS_URL, "test_url")
    m_js = Mock(spec=JetStreamContext)
    nc.jetstream = Mock(return_value=m_js)

    runner = Runner(nc=nc)
    with pytest.raises(Exception):
        await runner.initialize()

    assert getattr(runner, "js", None) is None


async def test_runner_initialize_jetstream_ko():
    nc = AsyncMock(spec=NatsClient)
    nc.connect.return_value = None
    v.set(NATS_URL, "test_url")
    nc.jetstream = Mock(side_effect=Exception("test exception"))

    runner = Runner(nc=nc)
    with pytest.raises(Exception):
        await runner.initialize()

    assert runner.nc is nc
    assert getattr(runner, "js", None) is None


@pytest.mark.parametrize(
    "runner_type, runner_method",
    [(TriggerRunner, "trigger_runner"), (TaskRunner, "task_runner"), (ExitRunner, "exit_runner")],
)
def test_get_runner_ok(runner_type, runner_method, m_runner):
    result = getattr(m_runner, runner_method)()

    assert isinstance(result, runner_type)
