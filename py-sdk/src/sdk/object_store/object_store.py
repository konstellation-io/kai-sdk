import re
from dataclasses import dataclass
from typing import Optional

from exceptions import (
    EmptyPayloadError,
    FailedCompilingRegexpError,
    FailedDeletingFileError,
    FailedGettingFileError,
    FailedListingFilesError,
    FailedObjectStoreInitializationError,
    FailedPurgingFilesError,
    FailedSavingFileError,
    UndefinedObjectStoreError,
)
from loguru import logger
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from nats.js.object_store import NatsObjectStore
from vyper import v


@dataclass
class ObjectStore:
    jetstream: JetStreamContext
    object_store_name: Optional[str] = ""
    object_store: Optional[NatsObjectStore] = None

    logger = logger.bind(context="[OBJECT STORE]")

    def __post_init__(self):
        if setattr(self, self.object_store_name, v.get("nats.object_store")):
            self.logger = logger.bind(context=f"[OBJECT STORE: {self.object_store_name}]")

    async def initialize(self):
        if self.object_store_name:
            self.object_store = await self.init_object_store_deps()

    async def init_object_store_deps(self) -> Optional[NatsObjectStore]:
        if self.object_store_name:
            try:
                object_store = await self.jetstream.object_store(self.object_store_name)
                return object_store
            except Exception as e:
                raise FailedObjectStoreInitializationError(error=e)

        self.logger.info("object store not defined. Skipping object store initialization")

    async def list(self, regexp: str = None) -> list:
        if not self.object_store:
            raise UndefinedObjectStoreError

        try:
            obj_store_list = await self.object_store.list()
        except Exception as e:
            raise FailedListingFilesError(error=e)

        pattern = None
        if regexp:
            try:
                pattern = re.compile(regexp)
            except re.error as e:
                raise FailedCompilingRegexpError(error=e)

        response = []

        for obj_name in obj_store_list:
            if not pattern or pattern.match(obj_name.name):
                response.append(obj_name.name)

        self.logger.info(f"files successfully listed from object store {self.object_store_name}")

        return response

    async def get(self, key: str) -> bytes:
        if not self.object_store:
            raise UndefinedObjectStoreError

        try:
            response = await self.object_store.get_bytes(key)
        except Exception as e:
            raise FailedGettingFileError(key=key, error=e)

        self.logger.info(f"file {key} successfully retrieved from object store {self.object_store_name}")

        return response

    async def save(self, key: str, payload: bytes):
        if not self.object_store:
            raise UndefinedObjectStoreError

        if not payload:
            raise EmptyPayloadError

        try:
            await self.object_store.put_bytes(key, payload)
        except Exception as e:
            raise FailedSavingFileError(key=key, error=e)

        self.logger.info(f"file {key} successfully saved in object store {self.object_store_name}")

    async def delete(self, key: str):
        if not self.object_store:
            raise UndefinedObjectStoreError

        try:
            await self.object_store.delete(key)
        except Exception as e:
            raise FailedDeletingFileError(key=key, error=e)

        self.logger.info(f"file {key} successfully deleted from object store {self.object_store_name}")

    async def purge(self, regexp: str = None):
        if not self.object_store:
            raise UndefinedObjectStoreError

        pattern = None
        if regexp:
            try:
                pattern = re.compile(regexp)
            except re.error as e:
                raise FailedCompilingRegexpError(error=e)

        objects = await self.list()

        for object_name in objects:
            if not pattern or pattern.match(object_name):
                self.logger.info(f"deleting file {object_name} from object store {self.object_store_name}")

                try:
                    await self.object_store.delete(object_name)
                except Exception as e:
                    raise FailedPurgingFilesError(error=e)

        self.logger.info(f"files successfully purged from object store {self.object_store_name}")
