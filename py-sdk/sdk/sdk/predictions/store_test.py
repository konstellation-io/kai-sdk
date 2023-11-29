from dataclasses import asdict
from datetime import datetime
from unittest.mock import Mock, patch

import pytest
from redis import Redis
from vyper import v

from sdk.metadata.metadata import Metadata
from sdk.predictions.exceptions import (
    EmptyIdError,
    FailedToFindPredictionsError,
    FailedToGetPredictionError,
    FailedToInitializePredictionsStoreError,
    FailedToSavePredictionError,
    FailedToUpdatePredictionError,
    MalformedEndpointError,
    MissingRequiredFilterFieldError,
    NotFoundError,
)
from sdk.predictions.store import Predictions
from sdk.predictions.types import Filter, Prediction, TimestampRange

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


@pytest.fixture
def m_prediction():
    return Prediction(
        creation_date=datetime.now(),
        last_modified=datetime.now(),
        payload={"test": "test"},
        metadata={
            "version": "test_version",
            "workflow": "test_workflow",
            "workflow_type": "test_workflow_type",
            "process": "test_process",
            "request_id": "test_request_id",
        },
    )


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


def test_update_ko(m_store):
    m_store.client.json.return_value.get.side_effect = Exception

    with pytest.raises(FailedToSavePredictionError):
        m_store.update("test_id", lambda x: x)


def test_update_wrong_id(m_store):
    with pytest.raises(FailedToSavePredictionError) as error:
        m_store.update("", lambda x: x)

        assert error == FailedToSavePredictionError("", EmptyIdError())


def test_update_failed_to_save_prediction(m_store):
    expected_prediction = Prediction(
        creation_date=datetime.now(),
        last_modified=datetime.now(),
        payload={"test": "test"},
        metadata={
            "version": "test_version",
            "workflow": "test_workflow",
            "workflow_type": "test_workflow_type",
            "process": "test_process",
            "request_id": "test_request_id",
        },
    )
    m_store.client.json.return_value.get.return_value = asdict(expected_prediction)
    m_store.client.json.return_value.set.side_effect = Exception

    with pytest.raises(FailedToSavePredictionError):
        m_store.update("test_id", lambda x: x)


def test_validate_filter_ok(m_store):
    filter_ = Filter(
        creation_date=TimestampRange(start_date=0, end_date=1),
    )
    m_store._validate_filter(filter_)

    assert filter_.version == Metadata.get_version()


def test_validate_filter_missing_required_filter_field_ko(m_store):
    with pytest.raises(MissingRequiredFilterFieldError):
        m_store._validate_filter(
            Filter(
                creation_date=TimestampRange(start_date=None, end_date=1),
            )
        )


def test_build_query_ok(m_store):
    time_ = datetime.now()
    expected_time = int(time_.timestamp() * 1000)
    result = m_store._build_query(
        Filter(
            version="test",
            workflow="test",
            workflow_type="test",
            process="test",
            request_id="test",
            creation_date=TimestampRange(start_date=time_, end_date=time_),
        )
    )

    assert (
        result
        == "@product:{%s} @creation_date:[%s %s] @version:{%s} @workflow:{%s} @workflow_type:{%s} @process:{%s} @request_id:{%s}"
        % (Metadata.get_product(), expected_time, expected_time, "test", "test", "test", "test", "test")
    )


def test_build_query_optional_fields_ok(m_store):
    v.set("metadata.version", "test_version")
    time_ = datetime.now()
    expected_time = int(time_.timestamp() * 1000)
    result = m_store._build_query(
        Filter(
            creation_date=TimestampRange(start_date=time_, end_date=time_),
        )
    )

    assert result == "@product:{%s} @creation_date:[%s %s] @version:{%s}" % (
        Metadata.get_product(),
        expected_time,
        expected_time,
        None,  # version is not set by default here
    )
