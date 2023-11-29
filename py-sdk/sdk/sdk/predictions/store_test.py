from unittest.mock import Mock, patch

import pytest
from redis import Redis
from vyper import v

from sdk.predictions.exceptions import (
    FailedToFindPredictionsError,
    FailedToGetPredictionError,
    FailedToInitializePredictionsStoreError,
    FailedToParseResultError,
    FailedToSavePredictionError,
    NotFoundError,
)
from sdk.predictions.store import Predictions


@pytest.fixture
def m_redis():
    return Mock(spec=Redis)


@pytest.fixture
def m_store(m_redis):
    store = Predictions()
    store.client = m_redis

    return store


@patch.object(Redis, "__init__", return_value=None)
def test_ok(m_redis_init, m_redis):
    v.set("predictions.endpoint", "test_endpoint")
    v.set("predictions.username", "test_username")
    v.set("predictions.password", "test_password")

    store = Predictions()

    m_redis_init.assert_called_once_with(host="test_endpoint", username="test_username", password="test_password")
    assert store.client is not None


@patch.object(Redis, "__init__", side_effect=Exception)
def test_ko(m_redis_init):
    v.set("predictions.endpoint", "test_endpoint")
    v.set("predictions.username", "test_username")
    v.set("predictions.password", "test_password")

    with pytest.raises(FailedToInitializePredictionsStoreError):
        Predictions()
