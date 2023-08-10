from typing import List
from unittest.mock import call
from mock import AsyncMock

import pytest
from nats.aio.client import Client as NatsClient
from nats.js.api import ObjectInfo
from nats.js.client import JetStreamContext
from nats.js.errors import NotFoundError, ObjectNotFoundError
from nats.js.object_store import ObjectStore as NatsObjectStore
from object_store.exceptions import (
    FailedCompilingRegexpError,
    FailedDeletingFileError,
    FailedGettingFileError,
    FailedListingFilesError,
    FailedObjectStoreInitializationError,
    FailedPurgingFilesError,
    FailedSavingFileError,
    UndefinedObjectStoreError,
)
from object_store.object_store import ObjectStore

KEY_140 = "object:140"
KEY_141 = "object:141"
KEY_142 = "object:142"
LIST_KEYS = [KEY_140, KEY_141, KEY_142]


@pytest.fixture(scope="function")
def m_objects() -> List[ObjectInfo]:
    objects = []
    for key in LIST_KEYS:
        object_info = ObjectInfo(
            name=key,
            deleted=False,
            bucket="test_bucket",
            nuid="test_nuid",
        )
        objects.append(object_info)

    return objects


@pytest.fixture(scope="function")
def m_object_store() -> ObjectStore:
    js = AsyncMock(spec=JetStreamContext)

    object_store = ObjectStore(js=js, object_store_name="test_object_store")
    object_store.object_store = AsyncMock(spec=NatsObjectStore)

    return object_store


def test_ok():
    nc = NatsClient()
    js = nc.jetstream()
    name = "test_object_store"

    object_store = ObjectStore(js=js, object_store_name=name)

    assert object_store.js is not None
    assert object_store.object_store_name == name
    assert object_store.object_store is None


async def test_initialize_ok(m_object_store):
    m_object_store.object_store = None
    fake_object_store = AsyncMock(spec=ObjectStore)
    m_object_store.js.object_store.return_value = fake_object_store

    await m_object_store.initialize()

    assert m_object_store.object_store == fake_object_store


async def test_initialize_undefined_object_store_name(m_object_store):
    m_object_store.object_store_name = None
    m_object_store.object_store = None

    await m_object_store.initialize()

    assert m_object_store.object_store is None


async def test_initialize_ko(m_object_store):
    m_object_store.object_store = None
    m_object_store.js.object_store.side_effect = Exception

    with pytest.raises(FailedObjectStoreInitializationError):
        await m_object_store.initialize()


async def test_list_ok(m_object_store, m_objects):
    m_object_store.object_store.list.return_value = m_objects

    result = await m_object_store.list()

    assert m_object_store.object_store.list.called
    assert result == [obj.name for obj in m_objects]


async def test_list_regex_ok(m_object_store, m_objects):
    m_object_store.object_store.list.return_value = m_objects
    expected = [m_objects[0].name, m_objects[1].name]

    result = await m_object_store.list(r"(object:140|object:141)")

    assert m_object_store.object_store.list.called
    assert result == expected


async def test_list_regex_ko(m_object_store, m_objects):
    m_object_store.object_store.list.return_value = m_objects

    with pytest.raises(FailedCompilingRegexpError):
        await m_object_store.list(1)


async def test_list_undefined_ko(m_object_store):
    m_object_store.object_store_name = None
    m_object_store.object_store = None

    with pytest.raises(UndefinedObjectStoreError):
        await m_object_store.list()


async def test_list_not_found(m_object_store):
    m_object_store.object_store.list.side_effect = NotFoundError

    result = await m_object_store.list()

    assert m_object_store.object_store.list.called
    assert result == []


async def test_list_failed_ko(m_object_store):
    m_object_store.object_store.list.side_effect = Exception
    with pytest.raises(FailedListingFilesError):
        await m_object_store.list()


async def test_get_ok(m_object_store):
    expected = AsyncMock(spec=ObjectInfo)
    m_object_store.object_store.get.return_value = expected

    result = await m_object_store.get("test-key")

    assert m_object_store.object_store.get.called
    assert m_object_store.object_store.get.call_args == call("test-key")
    assert result == (expected, True)


async def test_get_undefined_ko(m_object_store):
    m_object_store.object_store_name = None
    m_object_store.object_store = None

    with pytest.raises(UndefinedObjectStoreError):
        await m_object_store.get("key-1")


async def test_get_not_found(m_object_store):
    m_object_store.object_store.get.side_effect = ObjectNotFoundError

    result = await m_object_store.get("test-key")

    assert m_object_store.object_store.get.called
    assert result == (None, False)


async def test_get_failed_ko(m_object_store):
    m_object_store.object_store.get.side_effect = Exception

    with pytest.raises(FailedGettingFileError):
        await m_object_store.get("test-key")


async def test_save_ok(m_object_store):
    result = await m_object_store.save("test-key", b"any")

    assert m_object_store.object_store.put.called
    assert m_object_store.object_store.put.call_args == call("test-key", b"any")
    assert result is None


async def test_save_undefined_ko(m_object_store):
    m_object_store.object_store_name = None
    m_object_store.object_store = None

    with pytest.raises(UndefinedObjectStoreError):
        await m_object_store.save("key-1", b"any2")


async def test_save_missing_payload_ko(m_object_store):
    with pytest.raises(Exception):
        await m_object_store.save("key-1")


async def test_save_failed_ko(m_object_store):
    m_object_store.object_store.put.side_effect = Exception

    with pytest.raises(FailedSavingFileError):
        await m_object_store.save("test-key", b"prueba")


async def test_delete_ok(m_object_store, m_objects):
    m_object_store.object_store.list.return_value = m_objects
    deleted_object = m_objects[0]
    deleted_object.deleted = True
    m_object_store.object_store.delete.return_value=deleted_object

    result = await m_object_store.delete(KEY_140)

    assert result


async def test_delete_undefined_ko(m_object_store):
    m_object_store.object_store_name = None
    m_object_store.object_store = None

    with pytest.raises(UndefinedObjectStoreError):
        await m_object_store.delete("key-1")


async def test_delete_not_found_ko(m_object_store):
    m_object_store.object_store.delete.side_effect=ObjectNotFoundError

    result = await m_object_store.delete("key-1")

    assert not result


async def test_delete_failed_ko(m_object_store):
    m_object_store.object_store.delete.side_effect = Exception

    with pytest.raises(FailedDeletingFileError):
        await m_object_store.delete("test-key")


async def test_purge_ok(m_object_store, m_objects):
    m_object_store.object_store.list.return_value = [obj for obj in m_objects]
    for obj in m_objects:
        obj.deleted = True
    m_object_store.object_store.delete.side_effect=m_objects

    result = await m_object_store.purge()

    assert result is None
    assert m_object_store.object_store.delete.call_count == 3
    assert m_object_store.object_store.delete.call_args_list == [
        call(KEY_140),
        call(KEY_141),
        call(KEY_142),
    ]


async def test_purge_regex_ok(m_object_store, m_objects):
    m_object_store.object_store.list.return_value = [m_objects[0], m_objects[2]]
    for obj in m_objects:
        obj.deleted = True
    m_object_store.object_store.delete.side_effect=[m_objects[0], m_objects[2]]

    result = await m_object_store.purge(r"(object:140|object:142)")

    assert result is None
    assert m_object_store.object_store.delete.call_count == 2
    assert m_object_store.object_store.delete.call_args_list == [call(KEY_140), call(KEY_142)]


async def test_purge_undefined_ko(m_object_store):
    m_object_store.object_store_name = None
    m_object_store.object_store = None

    with pytest.raises(UndefinedObjectStoreError):
        await m_object_store.purge()


async def test_purge_regex_ko(m_object_store):
    with pytest.raises(FailedCompilingRegexpError):
        await m_object_store.purge(1)


async def test_purge_failed_ko(m_object_store, m_objects):
    m_object_store.object_store.list.return_value = [m_objects[0], m_objects[2]]
    m_object_store.object_store.delete.side_effect = Exception

    with pytest.raises(FailedPurgingFilesError):
        await m_object_store.purge()
