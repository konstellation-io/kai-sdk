import re
from dataclasses import dataclass
from typing import Optional

from loguru import logger
from loguru._logger import Logger
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
from vyper import v

UNDEFINED_OBJECT_STORE = "object store not defined"


@dataclass
class ObjectStore:
    js: JetStreamContext
    object_store_name: Optional[str] = None
    object_store: Optional[NatsObjectStore] = None
    logger: Logger = logger.bind(context="[OBJECT STORE]")

    def __post_init__(self):
        self.object_store_name = v.get("nats.object_store")
        if self.object_store_name:
            self.logger = logger.bind(context=f"[OBJECT STORE: {self.object_store_name}]")

    async def initialize(self) -> Optional[Exception]:
        if self.object_store_name:
            self.object_store = await self._init_object_store()

    async def _init_object_store(self) -> Optional[NatsObjectStore] | Exception:
        if self.object_store_name:
            try:
                object_store = await self.js.object_store(self.object_store_name)
                return object_store
            except Exception as e:
                logger.warning(f"failed initializing object store {self.object_store_name}: {e}")
                raise FailedObjectStoreInitializationError(error=e)

        self.logger.info("object store not defined [skipped]")

    async def list(self, regexp: Optional[str] = None) -> list[str] | Exception:
        if not self.object_store:
            logger.warning(UNDEFINED_OBJECT_STORE)
            raise UndefinedObjectStoreError

        try:
            objects = await self.object_store.list(ignore_deleted=True)
        except NotFoundError as e:
            logger.debug(f"not found any files in object store {self.object_store_name}: {e}")
            return []
        except Exception as e:
            logger.warning(f"failed listing files from object store {self.object_store_name}: {e}")
            raise FailedListingFilesError(error=e)

        pattern = None
        if regexp:
            try:
                pattern = re.compile(regexp)
            except Exception as e:
                logger.warning(f"failed compiling regexp {regexp}: {e}")
                raise FailedCompilingRegexpError(error=e)

        response = []
        for obj_name in objects:
            if not pattern or pattern.match(obj_name.name):
                response.append(obj_name)

        self.logger.info(f"files successfully listed from object store {self.object_store_name}")

        return response

    async def get(self, key: str) -> tuple[bytes, bool] | Exception:
        if not self.object_store:
            logger.warning(UNDEFINED_OBJECT_STORE)
            raise UndefinedObjectStoreError

        try:
            response = await self.object_store.get(key)
        except ObjectNotFoundError:
            logger.debug(f"file {key} not found in object store {self.object_store_name}")
            return None, False
        except Exception as e:
            logger.warning(f"failed getting file {key} from object store {self.object_store_name}: {e}")
            raise FailedGettingFileError(key=key, error=e)
        else:
            self.logger.info(f"file {key} successfully retrieved from object store {self.object_store_name}")

        return response, True

    async def save(self, key: str, payload: bytes) -> Optional[Exception]:
        if not self.object_store:
            logger.warning(UNDEFINED_OBJECT_STORE)
            raise UndefinedObjectStoreError

        try:
            await self.object_store.put(key, payload)
        except Exception as e:
            logger.warning(f"failed saving file {key} in object store {self.object_store_name}: {e}")
            raise FailedSavingFileError(key=key, error=e)
        else:
            self.logger.info(f"file {key} successfully saved in object store {self.object_store_name}")

    async def delete(self, key: str) -> bool | Exception:
        if not self.object_store:
            logger.warning(UNDEFINED_OBJECT_STORE)
            raise UndefinedObjectStoreError

        try:
            info_ = self.object_store.delete(key)
            return info_.deleted
        except ObjectNotFoundError:
            logger.debug(f"file {key} not found in object store {self.object_store_name}")
            return False
        except Exception as e:
            logger.warning(f"failed deleting file {key} from object store {self.object_store_name}: {e}")
            raise FailedDeletingFileError(key=key, error=e)

    async def purge(self, regexp: Optional[str] = None) -> Optional[Exception]:
        if not self.object_store:
            logger.warning(UNDEFINED_OBJECT_STORE)
            raise UndefinedObjectStoreError

        pattern = None
        if regexp:
            try:
                pattern = re.compile(regexp)
            except Exception as e:
                logger.warning(f"failed compiling regexp {regexp}: {e}")
                raise FailedCompilingRegexpError(error=e)

        objects = await self.list()
        deleted = 0
        for _, obj_name in objects:
            if not pattern or pattern.match(obj_name):
                self.logger.info(f"deleting file {obj_name} from object store {self.object_store_name}...")

                try:
                    if await self.object_store.delete(obj_name).info.deleted:
                        deleted += 1
                        logger.info(f"file {obj_name} successfully deleted from object store {self.object_store_name}")
                except Exception as e:
                    logger.warning(f"failed deleting file {obj_name} from object store {self.object_store_name}: {e}")
                    raise FailedPurgingFilesError(error=e)

        self.logger.info(f"{deleted} files successfully purged from object store {self.object_store_name}")
