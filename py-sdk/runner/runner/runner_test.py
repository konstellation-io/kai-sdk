from mock import patch
from nats.aio.client import Client as NatsClient
from vyper import v

from runner.kai_nats_msg_pb2 import KaiNatsMessage
from runner.runner import Runner
from sdk.centralized_config.centralized_config import CentralizedConfig
from sdk.kai_sdk import KaiSDK

PRODUCT_BUCKET = "centralized_configuration.product.bucket"
WORKFLOW_BUCKET = "centralized_configuration.workflow.bucket"
PROCESS_BUCKET = "centralized_configuration.process.bucket"
NATS_OBJECT_STORE = "nats.object_store"


@patch.object(CentralizedConfig, "_init_kv_stores", return_value=("test_product", "test_workflow", "test_process"))
async def test_sdk_import_ok(_):
    nc = NatsClient()
    js = nc.jetstream()
    req_msg = KaiNatsMessage()
    v.set(NATS_OBJECT_STORE, None)
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    sdk = KaiSDK(nc=nc, js=js, req_msg=req_msg)
    await sdk.initialize()

    assert sdk.nc is not None
    assert sdk.js is not None
    assert sdk.req_msg is not None
    assert sdk.logger is not None
    assert sdk.metadata is not None
    assert sdk.messaging is not None
    assert sdk.object_store.object_store_name is None
    assert sdk.object_store.object_store is None
    assert sdk.centralized_config is not None
    assert sdk.centralized_config.product_kv == "test_product"
    assert sdk.centralized_config.workflow_kv == "test_workflow"
    assert sdk.centralized_config.process_kv == "test_process"
    assert sdk.path_utils is not None
    assert sdk.measurements is None
    assert sdk.storage is None


@patch.object(NatsClient, "connect")
async def test_runner_ok(_):
    nc = NatsClient()
    js = nc.jetstream()

    runner = Runner(nc=nc, js=js)
    # await runner.initialize()
