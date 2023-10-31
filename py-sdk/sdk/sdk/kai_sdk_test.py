import pytest
from mock import AsyncMock, Mock, patch
from nats.aio.client import Client as NatsClient
from nats.js.kv import KeyValue
from nats.js.object_store import ObjectStore as NatsObjectStore
from vyper import v

from sdk.centralized_config.centralized_config import CentralizedConfig
from sdk.centralized_config.exceptions import FailedToInitializeConfigError
from sdk.ephemeral_storage.ephemeral_storage import EphemeralStorage
from sdk.ephemeral_storage.exceptions import FailedToInitializeEphemeralStorageError
from sdk.kai_nats_msg_pb2 import KaiNatsMessage
from sdk.kai_sdk import KaiSDK, MeasurementsABC, Storage
from sdk.messaging.messaging import Messaging
from sdk.metadata.metadata import Metadata
from sdk.path_utils.path_utils import PathUtils
from sdk.persistent_storage.persistent_storage import PersistentStorage

GLOBAL_BUCKET = "centralized_configuration.global.bucket"
PRODUCT_BUCKET = "centralized_configuration.product.bucket"
WORKFLOW_BUCKET = "centralized_configuration.workflow.bucket"
PROCESS_BUCKET = "centralized_configuration.process.bucket"
NATS_OBJECT_STORE = "nats.object_store"


@patch.object(
    CentralizedConfig,
    "_init_kv_stores",
    return_value=(
        AsyncMock(spec=KeyValue),
        AsyncMock(spec=KeyValue),
        AsyncMock(spec=KeyValue),
        AsyncMock(spec=KeyValue),
    ),
)
@patch.object(PersistentStorage, "__new__", return_value=Mock(spec=PersistentStorage))
async def test_initialize_ok(persistent_storage_mock, centralized_config_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()
    v.set(NATS_OBJECT_STORE, None)
    v.set(GLOBAL_BUCKET, "test_global")
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    sdk = KaiSDK(nc=nc, js=js)
    await sdk.initialize()

    assert isinstance(sdk.metadata, Metadata)
    assert isinstance(sdk.messaging, Messaging)
    assert isinstance(sdk.storage, Storage)
    assert isinstance(sdk.storage.ephemeral, EphemeralStorage)
    assert isinstance(sdk.storage.persistent, PersistentStorage)
    assert isinstance(sdk.centralized_config, CentralizedConfig)
    assert isinstance(sdk.path_utils, PathUtils)
    assert sdk.nc is not None
    assert sdk.js is not None
    assert getattr(sdk, "request_msg", None) is None
    assert sdk.logger is not None
    assert sdk.metadata is not None
    assert sdk.messaging is not None
    assert getattr(sdk.messaging, "request_msg", None) is None
    assert sdk.storage is not None
    assert sdk.storage.ephemeral is not None
    assert sdk.storage.ephemeral.ephemeral_storage_name == ""
    assert sdk.centralized_config is not None
    assert isinstance(sdk.centralized_config.global_kv, KeyValue)
    assert isinstance(sdk.centralized_config.product_kv, KeyValue)
    assert isinstance(sdk.centralized_config.workflow_kv, KeyValue)
    assert isinstance(sdk.centralized_config.process_kv, KeyValue)
    assert sdk.path_utils is not None
    assert isinstance(sdk.measurements, MeasurementsABC)


@patch.object(CentralizedConfig, "_init_kv_stores", side_effect=Exception)
@patch.object(PersistentStorage, "__new__", return_value=Mock(spec=PersistentStorage))
async def test_initialize_ko(persistent_storage_mock, centralized_config_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()
    v.set(NATS_OBJECT_STORE, None)
    v.set(GLOBAL_BUCKET, "test_global")
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    with pytest.raises(SystemExit):
        with pytest.raises(FailedToInitializeConfigError):
            sdk = KaiSDK(nc=nc, js=js)
            await sdk.initialize()


@patch.object(EphemeralStorage, "_init_object_store", return_value=Mock(spec=NatsObjectStore))
@patch.object(
    CentralizedConfig,
    "_init_kv_stores",
    return_value=(
        AsyncMock(spec=KeyValue),
        AsyncMock(spec=KeyValue),
        AsyncMock(spec=KeyValue),
        AsyncMock(spec=KeyValue),
    ),
)
@patch.object(PersistentStorage, "__new__", return_value=Mock(spec=PersistentStorage))
async def test_nats_initialize_ok(
    persistent_storage_mock, centralized_config_initialize_mock, object_store_initialize_mock
):
    nc = NatsClient()
    js = nc.jetstream()
    v.set(NATS_OBJECT_STORE, "test_object_store")
    v.set(GLOBAL_BUCKET, "test_global")
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    sdk = KaiSDK(nc=nc, js=js)
    await sdk.initialize()

    assert isinstance(sdk.centralized_config, CentralizedConfig)
    assert isinstance(sdk.storage, Storage)
    assert isinstance(sdk.storage.persistent, PersistentStorage)
    assert isinstance(sdk.storage.ephemeral, EphemeralStorage)
    assert sdk.centralized_config is not None
    assert isinstance(sdk.centralized_config.global_kv, KeyValue)
    assert isinstance(sdk.centralized_config.product_kv, KeyValue)
    assert isinstance(sdk.centralized_config.workflow_kv, KeyValue)
    assert isinstance(sdk.centralized_config.process_kv, KeyValue)
    assert sdk.storage is not None
    assert sdk.storage.ephemeral.object_store is not None
    assert sdk.storage.ephemeral.ephemeral_storage_name == "test_object_store"


@patch.object(EphemeralStorage, "_init_object_store", side_effect=Exception)
@patch.object(PersistentStorage, "__new__", return_value=Mock(spec=PersistentStorage))
async def test_nats_initialize_ko(_, object_store_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()

    with pytest.raises(SystemExit):
        with pytest.raises(FailedToInitializeEphemeralStorageError):
            sdk = KaiSDK(nc=nc, js=js)
            await sdk.initialize()


@patch.object(
    CentralizedConfig,
    "_init_kv_stores",
    return_value=(
        AsyncMock(spec=KeyValue),
        AsyncMock(spec=KeyValue),
        AsyncMock(spec=KeyValue),
        AsyncMock(spec=KeyValue),
    ),
)
@patch.object(PersistentStorage, "__new__", return_value=Mock(spec=PersistentStorage))
async def test_get_request_id_ok(persistent_storage_mock, centralized_config_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()
    request_msg = KaiNatsMessage(request_id="test_request_id")
    v.set(NATS_OBJECT_STORE, None)
    v.set(GLOBAL_BUCKET, "test_global")
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
    return_value=(
        AsyncMock(spec=KeyValue),
        AsyncMock(spec=KeyValue),
        AsyncMock(spec=KeyValue),
        AsyncMock(spec=KeyValue),
    ),
)
@patch.object(PersistentStorage, "__new__", return_value=Mock(spec=PersistentStorage))
async def test_set_request_msg_ok(persistent_storage_mock, centralized_config_initialize_mock):
    nc = NatsClient()
    js = nc.jetstream()
    request_msg = KaiNatsMessage(request_id="test_request_id")
    v.set(NATS_OBJECT_STORE, None)
    v.set(GLOBAL_BUCKET, "test_global")
    v.set(PRODUCT_BUCKET, "test_product")
    v.set(WORKFLOW_BUCKET, "test_workflow")
    v.set(PROCESS_BUCKET, "test_process")

    sdk = KaiSDK(nc=nc, js=js)
    await sdk.initialize()
    sdk.set_request_msg(request_msg)

    assert sdk.request_msg == request_msg
    assert isinstance(sdk.messaging, Messaging)
    assert sdk.messaging.request_msg == request_msg
