import pytest
from centralized_config.centralized_config import CentralizedConfig
from centralized_config.exceptions import FailedInitializingConfigError
from mock import Mock
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from nats.js.kv import KeyValue
from vyper import v


@pytest.fixture(scope="function")
def m_centralized_config() -> CentralizedConfig:
    js = Mock(spec=JetStreamContext)
    product_kv = Mock(spec=KeyValue)
    workflow_kv = Mock(spec=KeyValue)
    process_kv = Mock(spec=KeyValue)

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
    fake_product_kv = Mock(spec=KeyValue)
    fake_workflow_kv = Mock(spec=KeyValue)
    fake_process_kv = Mock(spec=KeyValue)
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
