import pytest
from mock import AsyncMock, patch
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v

from runner.kai_nats_msg_pb2 import KaiNatsMessage
from runner.runner import Runner
from sdk.centralized_config.centralized_config import CentralizedConfig
from sdk.kai_sdk import KaiSDK

PRODUCT_BUCKET = "centralized_configuration.product.bucket"
WORKFLOW_BUCKET = "centralized_configuration.workflow.bucket"
PROCESS_BUCKET = "centralized_configuration.process.bucket"
NATS_OBJECT_STORE = "nats.object_store"


@pytest.fixture(scope="function")
def m_runner() -> Runner:
    nc = AsyncMock(spec=NatsClient)
    js = AsyncMock(spec=JetStreamContext)
    v.set("nats.url", "test_url")
    v.set("APP_CONFIG_PATH", "test_path")
    v.set("runner.logger.output_paths", ["stdout"])
    v.set("runner.logger.error_output_paths", ["stderr"])
    v.set("runner.logger.level", "INFO")

    runner = Runner(nc=nc, js=js)

    return runner


@patch.object(CentralizedConfig, "_init_kv_stores", return_value=("test_product", "test_workflow", "test_process"))
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

    assert sdk.nc is not None
    assert sdk.js is not None
    assert sdk.req_msg == req_msg
    assert sdk.logger is not None  # TODO add logger initializing in runner test
    assert sdk.metadata is not None
    assert sdk.messaging is not None
    assert sdk.messaging.req_msg == req_msg
    assert sdk.object_store.object_store_name is None
    assert sdk.object_store.object_store is None
    assert sdk.centralized_config is not None
    assert sdk.centralized_config.product_kv == "test_product"
    assert sdk.centralized_config.workflow_kv == "test_workflow"
    assert sdk.centralized_config.process_kv == "test_process"
    assert sdk.path_utils is not None
    assert sdk.measurements is None
    assert sdk.storage is None


@patch.object(NatsClient, "connect", return_value=None)
@patch("runner.runner.NatsClient.jetstream", return_value=AsyncMock(spec=JetStreamContext))
async def test_runner_ok(nats_jetstream, nats_connect):
    nc = NatsClient()
    v.set("nats.url", "test_url")
    v.set("APP_CONFIG_PATH", "test_path")
    v.set("runner.logger.output_paths", ["stdout"])
    v.set("runner.logger.error_output_paths", ["stderr"])
    v.set("runner.logger.level", "INFO")

    runner = Runner(nc=nc)

    assert runner.nc is not None
    assert runner.js is None
    assert runner.logger is not None
    # assert runner.logger.level == "INFO"
    # assert runner.logger._handlers[0].sink == "stdout"
    # assert runner.logger._handlers[1].sink == "stderr"


# @patch.object(NatsClient, "connect", return_value=None)
# @patch("runner.runner.NatsClient.jetstream", return_value=AsyncMock(spec=JetStreamContext))
# async def test_runner_ok(nats_jetstream, nats_connect):
#     nc = NatsClient()
#     v.set("nats.url", "test_url")

#     runner = Runner(nc=nc)
#     await runner.initialize()

#     assert runner.nc is not None
#     assert runner.js is not None
