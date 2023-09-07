import pytest
from mock import AsyncMock, Mock, patch
from nats.aio.client import Client as NatsClient
from nats.js.kv import KeyValue
from nats.js.object_store import ObjectStore as NatsObjectStore
from vyper import v

from sdk.centralized_config.centralized_config import CentralizedConfig
from sdk.centralized_config.exceptions import FailedInitializingConfigError
from sdk.kai_nats_msg_pb2 import KaiNatsMessage
from sdk.kai_sdk import KaiSDK, MeasurementsABC, StorageABC
from sdk.messaging.messaging import Messaging
from sdk.metadata.metadata import Metadata
from sdk.object_store.object_store import FailedObjectStoreInitializationError, ObjectStore
from sdk.path_utils.path_utils import PathUtils

PRODUCT_BUCKET = "centralized_configuration.product.bucket"
WORKFLOW_BUCKET = "centralized_configuration.workflow.bucket"
PROCESS_BUCKET = "centralized_configuration.process.bucket"
NATS_OBJECT_STORE = "nats.object_store"


@patch.object(
    CentralizedConfig,
    "_init_kv_stores",
    return_value=(AsyncMock(spec=KeyValue), AsyncMock(spec=KeyValue), AsyncMock(spec=KeyValue)),
)
async def test_initialize_ok(centralized_config_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()
    v.set(NATS_OBJECT_STORE, None)
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    sdk = KaiSDK(nc=nc, js=js)
    await sdk.initialize()

    assert isinstance(sdk.metadata, Metadata)
    assert isinstance(sdk.messaging, Messaging)
    assert isinstance(sdk.object_store, ObjectStore)
    assert isinstance(sdk.centralized_config, CentralizedConfig)
    assert isinstance(sdk.path_utils, PathUtils)
    assert sdk.nc is not None
    assert sdk.js is not None
    assert getattr(sdk, "request_msg", None) is None
    assert sdk.logger is not None
    assert sdk.metadata is not None
    assert sdk.messaging is not None
    assert getattr(sdk.messaging, "request_msg", None) is None
    assert sdk.object_store is not None
    assert sdk.object_store.object_store_name == ""
    assert sdk.object_store.object_store is None
    assert sdk.centralized_config is not None
    assert isinstance(sdk.centralized_config.process_kv, KeyValue)
    assert isinstance(sdk.centralized_config.product_kv, KeyValue)
    assert isinstance(sdk.centralized_config.workflow_kv, KeyValue)
    assert sdk.path_utils is not None
    assert isinstance(sdk.measurements, MeasurementsABC)
    assert isinstance(sdk.storage, StorageABC)


@patch.object(CentralizedConfig, "_init_kv_stores", side_effect=Exception)
async def test_initialize_ko(centralized_config_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()
    v.set(NATS_OBJECT_STORE, None)
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    with pytest.raises(SystemExit):
        with pytest.raises(FailedInitializingConfigError):
            sdk = KaiSDK(nc=nc, js=js)
            await sdk.initialize()


@patch.object(ObjectStore, "_init_object_store", return_value=Mock(spec=NatsObjectStore))
@patch.object(
    CentralizedConfig,
    "_init_kv_stores",
    return_value=(AsyncMock(spec=KeyValue), AsyncMock(spec=KeyValue), AsyncMock(spec=KeyValue)),
)
async def test_nats_initialize_ok(centralized_config_initialize_mock, object_store_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()
    v.set(NATS_OBJECT_STORE, "test_object_store")
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    sdk = KaiSDK(nc=nc, js=js)
    await sdk.initialize()

    assert isinstance(sdk.centralized_config, CentralizedConfig)
    assert isinstance(sdk.object_store, ObjectStore)
    assert sdk.centralized_config is not None
    assert isinstance(sdk.centralized_config.process_kv, KeyValue)
    assert isinstance(sdk.centralized_config.product_kv, KeyValue)
    assert isinstance(sdk.centralized_config.workflow_kv, KeyValue)
    assert sdk.object_store.object_store is not None
    assert sdk.object_store.object_store_name == "test_object_store"


@patch.object(ObjectStore, "_init_object_store", side_effect=Exception)
async def test_nats_initialize_ko(object_store_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()

    with pytest.raises(SystemExit):
        with pytest.raises(FailedObjectStoreInitializationError):
            sdk = KaiSDK(nc=nc, js=js)
            await sdk.initialize()


@patch.object(
    CentralizedConfig,
    "_init_kv_stores",
    return_value=(AsyncMock(spec=KeyValue), AsyncMock(spec=KeyValue), AsyncMock(spec=KeyValue)),
)
async def test_get_request_id_ok(centralized_config_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()
    request_msg = KaiNatsMessage(request_id="test_request_id")
    v.set(NATS_OBJECT_STORE, None)
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    sdk = KaiSDK(nc=nc, js=js)
    await sdk.initialize()

    assert sdk.get_request_id() is None

    sdk.set_request_msg(request_msg)

    assert sdk.get_request_id() == "test_request_id"


@patch.object(
    CentralizedConfig,
    "_init_kv_stores",
    return_value=(AsyncMock(spec=KeyValue), AsyncMock(spec=KeyValue), AsyncMock(spec=KeyValue)),
)
async def test_set_request_msg_ok(centralized_config_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()
    request_msg = KaiNatsMessage(request_id="test_request_id")
    v.set(NATS_OBJECT_STORE, None)
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    sdk = KaiSDK(nc=nc, js=js)
    await sdk.initialize()
    sdk.set_request_msg(request_msg)

    assert sdk.request_msg == request_msg
    assert isinstance(sdk.messaging, Messaging)
    assert sdk.messaging.request_msg == request_msg
