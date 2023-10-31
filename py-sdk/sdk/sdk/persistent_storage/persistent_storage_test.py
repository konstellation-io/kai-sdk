import io
from unittest.mock import Mock, patch

import pytest
import urllib3
from minio import Minio
from vyper import v

from sdk.persistent_storage.exceptions import (
    FailedToDeleteFileError,
    FailedToGetFileError,
    FailedToInitializePersistentStorageError,
    FailedToSaveFileError,
)
from sdk.persistent_storage.persistent_storage import PersistentStorage, PersistentStorageABC

TTL_DAYS = 30


@pytest.fixture(scope="function")
@patch.object(Minio, "__new__", return_value=Mock(spec=Minio))
def m_persistent_storage(minio_mock: Mock) -> PersistentStorageABC:
    persistent_storage = PersistentStorage()
    persistent_storage.minio_client = minio_mock
    persistent_storage.minio_bucket_name = "test-minio-bucket"

    return persistent_storage


@pytest.fixture(scope="function")
def m_object() -> urllib3.BaseHTTPResponse:
    object_ = Mock(spec=urllib3.BaseHTTPResponse)
    object_.close.return_value = None
    object_.read.return_value = b"test-payload"

    return object_


@patch.object(Minio, "__new__", return_value=Mock(spec=Minio))
def test_ok(_):
    v.set("minio.endpoint", "test-endpoint")
    v.set("minio.access_key_id", "test-access-key")
    v.set("minio.access_key_secret", "test-secret-key")
    v.set("minio.use_ssl", False)
    v.set("minio.bucket", "test-minio-bucket")

    persistent_storage = PersistentStorage()

    assert persistent_storage.minio_client is not None
    assert persistent_storage.minio_bucket_name == "test-minio-bucket"


def test_ko(m_persistent_storage):
    with pytest.raises(FailedToInitializePersistentStorageError):
        PersistentStorage()

        assert m_persistent_storage.minio_client is None


def test_save_ok(m_persistent_storage):
    m_persistent_storage.minio_client.put_object.return_value = None
    payload = io.BytesIO(b"test-payload")

    m_persistent_storage.save("test-key", payload, TTL_DAYS)

    m_persistent_storage.minio_client.put_object.assert_called_once()


def test_save_no_ttl_ok(m_persistent_storage):
    m_persistent_storage.minio_client.put_object.return_value = None
    payload = io.BytesIO(b"test-payload")

    m_persistent_storage.save("test-key", payload)

    m_persistent_storage.minio_client.put_object.assert_called_once_with(
        "test-minio-bucket",
        "test-key",
        payload,
        payload.getbuffer().nbytes,
    )


def test_save_ko(m_persistent_storage):
    m_persistent_storage.minio_client.put_object.side_effect = Exception
    payload = io.BytesIO(b"test-payload")

    with pytest.raises(FailedToSaveFileError):
        m_persistent_storage.save("test-key", payload, TTL_DAYS)

    m_persistent_storage.minio_client.put_object.assert_called_once()


@patch("sdk.persistent_storage.persistent_storage.PersistentStorage._object_exist", return_value=True)
def test_get_ok(_, m_persistent_storage, m_object):
    m_persistent_storage.minio_client.get_object.return_value = m_object

    payload, ok = m_persistent_storage.get("test-key", "test-version")

    m_persistent_storage.minio_client.get_object.assert_called_once()
    assert payload == m_object.read.return_value
    assert ok


@patch("sdk.persistent_storage.persistent_storage.PersistentStorage._object_exist", return_value=False)
def test_get_not_found_ok(_, m_persistent_storage):
    payload, ok = m_persistent_storage.get("test-key", "test-version")

    m_persistent_storage.minio_client.get_object.assert_not_called()
    assert payload is None
    assert not ok


@patch("sdk.persistent_storage.persistent_storage.PersistentStorage._object_exist", return_value=True)
def test_get_ko(_, m_persistent_storage):
    m_persistent_storage.minio_client.get_object.side_effect = Exception

    with pytest.raises(FailedToGetFileError):
        m_persistent_storage.get("test-key", "test-version")

    m_persistent_storage.minio_client.get_object.assert_called_once()


@patch("sdk.persistent_storage.persistent_storage.PersistentStorage._object_exist", return_value=True)
def test_delete_ok(_, m_persistent_storage):
    m_persistent_storage.minio_client.remove_object.return_value = None

    m_persistent_storage.delete("test-key")

    m_persistent_storage.minio_client.remove_object.assert_called_once()


@patch("sdk.persistent_storage.persistent_storage.PersistentStorage._object_exist", return_value=False)
def test_delete_not_found_ok(_, m_persistent_storage):
    m_persistent_storage.minio_client.remove_object.return_value = None

    m_persistent_storage.delete("test-key")

    m_persistent_storage.minio_client.remove_object.assert_not_called()


@patch("sdk.persistent_storage.persistent_storage.PersistentStorage._object_exist", return_value=True)
def test_delete_ko(_, m_persistent_storage):
    m_persistent_storage.minio_client.remove_object.side_effect = Exception

    with pytest.raises(FailedToDeleteFileError):
        m_persistent_storage.delete("test-key")

    m_persistent_storage.minio_client.remove_object.assert_called_once()


def test_list_ok(m_persistent_storage):
    m_persistent_storage.minio_client.list_objects.return_value = [Mock(object_name="test-key")]

    keys = m_persistent_storage.list()

    m_persistent_storage.minio_client.list_objects.assert_called_once()
    assert keys == ["test-key"]


def test_list_ko(m_persistent_storage):
    m_persistent_storage.minio_client.list_objects.side_effect = Exception

    keys = m_persistent_storage.list()

    m_persistent_storage.minio_client.list_objects.assert_called_once()
    assert keys == []


def test_list_versions_ok(m_persistent_storage):
    m_persistent_storage.minio_client.list_objects.return_value = [Mock(object_name="test-key")]

    keys = m_persistent_storage.list_versions("test-key")

    m_persistent_storage.minio_client.list_objects.assert_called_once()
    assert keys == ["test-key"]


def test_list_versions_ko(m_persistent_storage):
    m_persistent_storage.minio_client.list_objects.side_effect = Exception

    keys = m_persistent_storage.list_versions("test-key")

    m_persistent_storage.minio_client.list_objects.assert_called_once()
    assert keys == []


def test__object_exist_ok(m_persistent_storage):
    m_persistent_storage.minio_client.stat_object.return_value = None

    exist = m_persistent_storage._object_exist("test-key", "test-version")

    m_persistent_storage.minio_client.stat_object.assert_called_once()
    assert exist


def test__object_exist_ko(m_persistent_storage):
    m_persistent_storage.minio_client.stat_object.side_effect = Exception

    with pytest.raises(Exception):
        m_persistent_storage._object_exist("test-key", "test-version")

        m_persistent_storage.minio_client.stat_object.assert_called_once()
