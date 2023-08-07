from centralized_config.centralized_config import CentralizedConfig
from kai_nats_msg import KaiNatsMessage
from kai_sdk import KaiSDK
from mock import patch
from nats.aio.client import Client as NatsClient
from object_store.object_store import ObjectStore
from vyper import v

PRODUCT_BUCKET = "centralized_configuration.product.bucket"
WORKFLOW_BUCKET = "centralized_configuration.workflow.bucket"
PROCESS_BUCKET = "centralized_configuration.process.bucket"
NATS_OBJECT_STORE = "nats.object_store"


@patch.object(CentralizedConfig, "_init_kv_stores", return_value=("test_product", "test_workflow", "test_process"))
async def test_initialize_ok(centralized_config_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()
    req_msg = KaiNatsMessage()
    v.set(NATS_OBJECT_STORE, None)
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    sdk = KaiSDK(nc=nc, js=js, req_msg=req_msg)
    await sdk.initialize()

    assert sdk.centralized_config is not None
    assert sdk.centralized_config.product_kv == "test_product"
    assert sdk.centralized_config.workflow_kv == "test_workflow"
    assert sdk.centralized_config.process_kv == "test_process"
    assert sdk.object_store.object_store_name == None
    assert sdk.object_store.object_store == None


@patch.object(CentralizedConfig, "_init_kv_stores", return_value=("test_product", "test_workflow", "test_process"))
@patch.object(ObjectStore, "_init_object_store", return_value="Not Empty")
async def test_nats_initialize_ok(object_store_initialize_mock, centralized_config_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()
    req_msg = KaiNatsMessage()
    v.set(NATS_OBJECT_STORE, "test_object_store")
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    sdk = KaiSDK(nc=nc, js=js, req_msg=req_msg)
    await sdk.initialize()

    assert sdk.centralized_config is not None
    assert sdk.centralized_config.product_kv == "test_product"
    assert sdk.centralized_config.workflow_kv == "test_workflow"
    assert sdk.centralized_config.process_kv == "test_process"
    assert sdk.object_store.object_store is not None
    assert sdk.object_store.object_store_name == "test_object_store"


@patch.object(CentralizedConfig, "_init_kv_stores", return_value=("test_product", "test_workflow", "test_process"))
async def test_get_request_id_ok(centralized_config_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()
    req_msg = KaiNatsMessage(request_id="test_request_id")
    v.set(NATS_OBJECT_STORE, None)
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    sdk = KaiSDK(nc=nc, js=js, req_msg=req_msg)
    await sdk.initialize()

    assert sdk.get_request_id() == "test_request_id"
