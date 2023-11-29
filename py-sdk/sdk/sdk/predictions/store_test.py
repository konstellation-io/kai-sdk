from unittest.mock import Mock, patch

import pytest
from redis import Redis
from vyper import v

from sdk.predictions.exceptions import FailedToInitializePredictionsStoreError, MalformedEndpointError
from sdk.predictions.store import Predictions

PREDICTIONS_ENDPOINT_KEY = "predictions.endpoint"
PREDICTIONS_USERNAME_KEY = "predictions.username"
PREDICTIONS_PASSWORD_KEY = "predictions.password"
PREDICTIONS_INDEX_KEY_KEY = "predictions.index_key"
PREDICTIONS_ENDPOINT = "localhost:6379"
PREDICTIONS_USERNAME = "test_username"
PREDICTIONS_PASSWORD = "test_password"
PREDICTIONS_INDEX_KEY = "test_index_key"


@pytest.fixture
def m_redis():
    return Mock(spec=Redis)


@pytest.fixture
def m_store(m_redis):
    store = Predictions()
    store.client = m_redis

    return store


@patch.object(Redis, "__init__", return_value=None)
def test_ok(m_redis_init):
    v.set(PREDICTIONS_ENDPOINT_KEY, PREDICTIONS_ENDPOINT)
    v.set(PREDICTIONS_USERNAME_KEY, PREDICTIONS_USERNAME)
    v.set(PREDICTIONS_PASSWORD_KEY, PREDICTIONS_PASSWORD)
    v.set(PREDICTIONS_INDEX_KEY_KEY, PREDICTIONS_INDEX_KEY)

    store = Predictions()

    m_redis_init.assert_called_once_with(
        host="localhost", port=6379, username="test_username", password="test_password"
    )
    assert store.client is not None


def test_malformed_endpoint_ko():
    v.set(PREDICTIONS_ENDPOINT_KEY, "localhost")
    v.set(PREDICTIONS_USERNAME_KEY, PREDICTIONS_USERNAME)
    v.set(PREDICTIONS_PASSWORD_KEY, PREDICTIONS_PASSWORD)
    v.set(PREDICTIONS_INDEX_KEY_KEY, PREDICTIONS_INDEX_KEY)

    with pytest.raises(FailedToInitializePredictionsStoreError) as error:
        Predictions()

        assert error == MalformedEndpointError("localhost", "localhost")


@patch.object(Redis, "__init__", side_effect=Exception)
def test_initialization_ko(_):
    v.set(PREDICTIONS_ENDPOINT_KEY, PREDICTIONS_ENDPOINT)
    v.set(PREDICTIONS_USERNAME_KEY, PREDICTIONS_USERNAME)
    v.set(PREDICTIONS_PASSWORD_KEY, PREDICTIONS_PASSWORD)
    v.set(PREDICTIONS_INDEX_KEY_KEY, PREDICTIONS_INDEX_KEY)

    with pytest.raises(FailedToInitializePredictionsStoreError):
        Predictions()
