from __future__ import annotations

import re
from dataclasses import dataclass, field
from typing import Optional

import loguru
from loguru import logger
from nats.js.client import JetStreamContext
from nats.js.errors import NotFoundError, ObjectNotFoundError
from nats.js.object_store import ObjectStore as NatsObjectStore
from vyper import v

from sdk.object_store.exceptions import (
    FailedCompilingRegexpError,
    FailedDeletingFileError,
    FailedGettingFileError,
    FailedListingFilesError,
    FailedObjectStoreInitializationError,
    FailedPurgingFilesError,
    FailedSavingFileError,
    UndefinedObjectStoreError,
)

UNDEFINED_OBJECT_STORE = "object store not defined"


@dataclass
class ObjectStore:
    js: JetStreamContext
    object_store_name: Optional[str] = None
    object_store: Optional[NatsObjectStore] = None
    logger: loguru.Logger = logger.bind(context="[OBJECT STORE]")

    def __post_init__(self) -> None:
        self.object_store_name = v.get("nats.object_store")
        if self.object_store_name:
            self.logger = logger.bind(context=f"[OBJECT STORE: {self.object_store_name}]")

    async def initialize(self) -> None:
        if self.object_store_name:
            object_store = await self._init_object_store()
            self.object_store = object_store
        else:
            self.logger.info("object store not defined [skipped]")

    async def _init_object_store(self) -> NatsObjectStore:
        try:
            assert isinstance(self.object_store_name, str)
            object_store = await self.js.object_store(self.object_store_name)
            self.logger.debug(f"object store {self.object_store_name} successfully initialized")
            return object_store
        except Exception as e:
            self.logger.warning(f"failed initializing object store {self.object_store_name}: {e}")
            raise FailedObjectStoreInitializationError(error=e)

    async def list(self, regexp: Optional[str] = None) -> list[str]:
        if not self.object_store:
            self.logger.warning(UNDEFINED_OBJECT_STORE)
            raise UndefinedObjectStoreError

        try:
            objects = await self.object_store.list(ignore_deletes=True)
        except NotFoundError as e:
            self.logger.debug(f"no files found in object store {self.object_store_name}: {e}")
            return []
        except Exception as e:
            self.logger.warning(f"failed listing files from object store {self.object_store_name}: {e}")
            raise FailedListingFilesError(error=e)

        pattern = None
        if regexp:
            try:
                pattern = re.compile(regexp)
            except Exception as e:
                self.logger.warning(f"failed compiling regexp {regexp}: {e}")
                raise FailedCompilingRegexpError(error=e)

        response = []
        for obj in objects:
            obj_name = obj.name
            if not pattern or pattern.match(obj_name):
                response.append(obj_name)

        self.logger.info(f"files successfully listed from object store {self.object_store_name}")
        return response

    async def get(self, key: str) -> tuple[Optional[bytes], bool]:
        if not self.object_store:
            self.logger.warning(UNDEFINED_OBJECT_STORE)
            raise UndefinedObjectStoreError

        try:
            response = await self.object_store.get(key)
            self.logger.info(f"file {key} successfully retrieved from object store {self.object_store_name}")
            return response.data, True
        except ObjectNotFoundError as e:
            self.logger.debug(f"file {key} not found in object store {self.object_store_name}: {e}")
            return None, False
        except Exception as e:
            self.logger.warning(f"failed getting file {key} from object store {self.object_store_name}: {e}")
            raise FailedGettingFileError(key=key, error=e)

    async def save(self, key: str, payload: bytes) -> None:
        if not self.object_store:
            self.logger.warning(UNDEFINED_OBJECT_STORE)
            raise UndefinedObjectStoreError

        try:
            await self.object_store.put(key, payload)
            self.logger.info(f"file {key} successfully saved in object store {self.object_store_name}")
        except Exception as e:
            self.logger.warning(f"failed saving file {key} in object store {self.object_store_name}: {e}")
            raise FailedSavingFileError(key=key, error=e)

    async def delete(self, key: str) -> bool:
        if not self.object_store:
            self.logger.warning(UNDEFINED_OBJECT_STORE)
            raise UndefinedObjectStoreError

        try:
            info_ = await self.object_store.delete(key)
            return info_.info.deleted if info_.info.deleted else False
        except ObjectNotFoundError as e:
            self.logger.debug(f"file {key} not found in object store {self.object_store_name}: {e}")
            return False
        except Exception as e:
            self.logger.warning(f"failed deleting file {key} from object store {self.object_store_name}: {e}")
            raise FailedDeletingFileError(key=key, error=e)

    async def purge(self, regexp: Optional[str] = None) -> None:
        if not self.object_store:
            self.logger.warning(UNDEFINED_OBJECT_STORE)
            raise UndefinedObjectStoreError

        pattern = None
        if regexp:
            try:
                pattern = re.compile(regexp)
            except Exception as e:
                self.logger.warning(f"failed compiling regexp {regexp}: {e}")
                raise FailedCompilingRegexpError(error=e)

        object_names = await self.list()
        deleted = 0
        for name in object_names:
            if not pattern or pattern.match(name):
                self.logger.info(f"deleting file {name} from object store {self.object_store_name}...")

                try:
                    info_ = await self.object_store.delete(name)
                    if info_.info.deleted:
                        deleted += 1
                        self.logger.info(f"file {name} successfully deleted from object store {self.object_store_name}")
                except Exception as e:
                    self.logger.warning(f"failed deleting file {name} from object store {self.object_store_name}: {e}")
                    raise FailedPurgingFilesError(error=e)

        self.logger.info(f"{deleted} files successfully purged from object store {self.object_store_name}")
