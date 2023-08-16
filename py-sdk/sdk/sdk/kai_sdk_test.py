import pytest
from loguru import logger
from mock import Mock, patch
from nats.aio.client import Client as NatsClient
from vyper import v

from sdk.centralized_config.centralized_config import CentralizedConfig
from sdk.centralized_config.exceptions import FailedInitializingConfigError
from sdk.kai_nats_msg_pb2 import KaiNatsMessage
from sdk.kai_sdk import KaiSDK
from sdk.object_store.object_store import FailedObjectStoreInitializationError, ObjectStore

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


@patch.object(CentralizedConfig, "_init_kv_stores", side_effect=Exception)
async def test_initialize_ko(centralized_config_initialize_mock, caplog):
    nc = NatsClient()
    js = nc.jetstream()
    req_msg = KaiNatsMessage()
    v.set(NATS_OBJECT_STORE, None)
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    with pytest.raises(SystemExit):
        with pytest.raises(FailedInitializingConfigError):
            logger.disable("kai_sdk")
            sdk = KaiSDK(nc=nc, js=js, req_msg=req_msg)
            await sdk.initialize()
            logger.enable("kai_sdk")


@patch.object(ObjectStore, "_init_object_store", return_value=Mock(spec=ObjectStore))
@patch.object(CentralizedConfig, "_init_kv_stores", return_value=("test_product", "test_workflow", "test_process"))
async def test_nats_initialize_ok(centralized_config_initialize_mock, object_store_initialize_mock):
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


@patch.object(ObjectStore, "_init_object_store", side_effect=Exception)
async def test_nats_initialize_ko(object_store_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()
    req_msg = KaiNatsMessage()

    with pytest.raises(SystemExit):
        with pytest.raises(FailedObjectStoreInitializationError):
            logger.disable("kai_sdk")
            sdk = KaiSDK(nc=nc, js=js, req_msg=req_msg)
            await sdk.initialize()
            logger.enable("kai_sdk")


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