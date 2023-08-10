import pytest
from centralized_config.centralized_config import CentralizedConfig
from centralized_config.exceptions import FailedInitializingConfigError
from mock import AsyncMock
from unittest.mock import call
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from nats.js.kv import KeyValue
from vyper import v
from nats.js.errors import KeyNotFoundError
from centralized_config.constants import Scope
from centralized_config.exceptions import (
    FailedDeletingConfigError,
    FailedGettingConfigError,
    FailedInitializingConfigError,
    FailedSettingConfigError,
)

@pytest.fixture(scope="function")
def m_centralized_config() -> CentralizedConfig:
    js = AsyncMock(spec=JetStreamContext)
    product_kv = AsyncMock(spec=KeyValue)
    workflow_kv = AsyncMock(spec=KeyValue)
    process_kv = AsyncMock(spec=KeyValue)

    centralized_config = CentralizedConfig(js=js)
    centralized_config.product_kv = product_kv
    centralized_config.workflow_kv = workflow_kv
    centralized_config.process_kv = process_kv

    return centralized_config


def test_ok():
    nc = NatsClient()
    js = nc.jetstream()

    centralized_config = CentralizedConfig(js=js)

    assert centralized_config.js is not None
    assert centralized_config.product_kv is None
    assert centralized_config.workflow_kv is None
    assert centralized_config.process_kv is None


async def test_initialize_ok(m_centralized_config):
    m_centralized_config.product_kv = None
    m_centralized_config.workflow_kv = None
    m_centralized_config.process_kv = None
    fake_product_kv = AsyncMock(spec=KeyValue)
    fake_workflow_kv = AsyncMock(spec=KeyValue)
    fake_process_kv = AsyncMock(spec=KeyValue)
    v.set("centralized_configuration.product.bucket", "test_product_bucket")
    v.set("centralized_configuration.workflow.bucket", "test_workflow_bucket")
    v.set("centralized_configuration.process.bucket", "test_process_bucket")
    m_centralized_config.js.key_value.side_effect = [fake_product_kv, fake_workflow_kv, fake_process_kv]

    await m_centralized_config.initialize()

    assert m_centralized_config.product_kv == fake_product_kv
    assert m_centralized_config.workflow_kv == fake_workflow_kv
    assert m_centralized_config.process_kv == fake_process_kv


async def test_initialize_ko(m_centralized_config):
    m_centralized_config.js.key_value.side_effect = Exception

    with pytest.raises(FailedInitializingConfigError):
        await m_centralized_config.initialize()


async def test_get_config_with_scope_ok(m_centralized_config):
    m_centralized_config._get_config_from_scope = AsyncMock(return_value="test_config")

    config, found = await m_centralized_config.get_config("test_key", Scope.WorkflowScope)

    assert m_centralized_config._get_config_from_scope.awaited
    assert m_centralized_config._get_config_from_scope.call_args == call("test_key", Scope.WorkflowScope)
    assert config == "test_config"
    assert found

async def test_get_config_without_scope_ok(m_centralized_config):
    m_centralized_config._get_config_from_scope = AsyncMock(side_effect=[KeyNotFoundError, "test_config"])

    config, found = await m_centralized_config.get_config("test_key")

    assert m_centralized_config._get_config_from_scope.awaited
    assert m_centralized_config._get_config_from_scope.call_args_list == [call("test_key", Scope.ProcessScope), call("test_key", Scope.WorkflowScope)]    
    assert config == "test_config"
    assert found

async def test_get_config_with_scope_not_found(m_centralized_config):
    m_centralized_config._get_config_from_scope = AsyncMock(side_effect=KeyNotFoundError)

    config, found = await m_centralized_config.get_config("test_key", Scope.ProductScope)

    assert m_centralized_config._get_config_from_scope.awaited
    assert m_centralized_config._get_config_from_scope.call_args == call("test_key", Scope.ProductScope) 
    assert config is None
    assert not found

async def test_get_config_with_scope_ko(m_centralized_config):
    m_centralized_config._get_config_from_scope = AsyncMock(side_effect=Exception)

    with pytest.raises(Exception):
        await m_centralized_config.get_config("test_key", scope=Scope.ProductScope)

async def test_get_config_without_scope_not_found(m_centralized_config):
    m_centralized_config._get_config_from_scope = AsyncMock(side_effect=[KeyNotFoundError, KeyNotFoundError, KeyNotFoundError])

    config, found = await m_centralized_config.get_config("test_key")

    assert m_centralized_config._get_config_from_scope.awaited
    assert m_centralized_config._get_config_from_scope.call_args_list == [call("test_key", Scope.ProcessScope), call("test_key", Scope.WorkflowScope), call("test_key", Scope.ProductScope)]
    assert config is None
    assert not found

async def test_get_config_without_scope_ko(m_centralized_config):
    m_centralized_config._get_config_from_scope = AsyncMock(side_effect=Exception)

    with pytest.raises(Exception):
        await m_centralized_config.get_config("test_key")

async def test_set_config_with_scope_ok(m_centralized_config):
    await m_centralized_config.set_config("test_key", b"test_value", Scope.ProductScope)

    assert m_centralized_config.product_kv.put.awaited
    assert m_centralized_config.product_kv.put.call_args == call("test_key", b"test_value")

async def test_set_config_without_scope_ok(m_centralized_config):
    await m_centralized_config.set_config("test_key", b"test_value")

    assert m_centralized_config.process_kv.put.awaited
    assert m_centralized_config.process_kv.put.call_args == call("test_key", b"test_value")

async def test_set_config_with_scope_ko(m_centralized_config):
    m_centralized_config.product_kv.put.side_effect=Exception

    with pytest.raises(FailedSettingConfigError):
        await m_centralized_config.set_config("test_key", b"test_value", Scope.ProductScope)

async def test_set_config_without_scope_ko(m_centralized_config):
    m_centralized_config.process_kv.put.side_effect=Exception

    with pytest.raises(FailedSettingConfigError):
        await m_centralized_config.set_config("test_key", b"test_value")

async def test_delete_config_with_scope_ok(m_centralized_config):
    await m_centralized_config.delete_config("test_key", Scope.ProductScope)

    assert m_centralized_config.product_kv.delete.awaited
    assert m_centralized_config.product_kv.delete.call_args == call("test_key")

async def test_delete_config_without_scope_ok(m_centralized_config):
    await m_centralized_config.delete_config("test_key")

    assert m_centralized_config.process_kv.delete.awaited
    assert m_centralized_config.process_kv.delete.call_args == call("test_key")

async def test_delete_config_with_scope_ko(m_centralized_config):
    m_centralized_config.product_kv.delete.side_effect=Exception

    with pytest.raises(FailedDeletingConfigError):
        await m_centralized_config.delete_config("test_key", Scope.ProductScope)

async def test_delete_config_without_scope_ko(m_centralized_config):
    m_centralized_config.process_kv.delete.side_effect=Exception

    with pytest.raises(FailedDeletingConfigError):
        await m_centralized_config.delete_config("test_key")

async def test__get_config_from_scope_ok(m_centralized_config):
    entry = AsyncMock(spec=KeyValue.Entry)
    entry.key = "test_key"
    entry.value = b"test_value"
    m_centralized_config.product_kv.get = AsyncMock(return_value=entry)

    result = await m_centralized_config._get_config_from_scope("test_key", Scope.ProductScope)

    assert result == "test_value"

async def test__get_scoped_config_with_scope_ok(m_centralized_config):
    result = m_centralized_config._get_scoped_config(Scope.ProductScope)

    assert result == m_centralized_config.product_kv

async def test__get_scoped_config_without_scope_ok(m_centralized_config):
    result = m_centralized_config._get_scoped_config()

    assert result == m_centralized_config.process_kv