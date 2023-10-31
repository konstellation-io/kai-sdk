from __future__ import annotations

from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from datetime import datetime, timedelta
from typing import BinaryIO, Optional

import loguru
from loguru import logger
from minio import Minio
from minio.retention import COMPLIANCE, Retention
from vyper import v

from sdk.persistent_storage.exceptions import (
    FailedToDeleteFileError,
    FailedToGetFileError,
    FailedToInitializePersistentStorageError,
    FailedToListFilesError,
    FailedToSaveFileError,
    MissingBucketError,
)


@dataclass
class PersistentStorageABC(ABC):
    @abstractmethod
    def save(self, key: str, payload: bytes, ttl_days: Optional[int]) -> str | None:
        pass

    @abstractmethod
    def get(self, key: str, version: Optional[str]) -> tuple[Optional[bytes], bool]:
        pass

    @abstractmethod
    def list(self) -> list[str]:
        pass

    @abstractmethod
    def list(self, key: str) -> list[str]:
        pass

    @abstractmethod
    def delete(self, key: str, version: Optional[str]) -> bool:
        pass


@dataclass
class PersistentStorage(PersistentStorageABC):
    logger: loguru.Logger = logger.bind(context="[PERSISTENT STORAGE]")
    minio_client: Minio = field(init=False)
    minio_bucket_name: str = field(init=False)

    def __post_init__(self) -> None:
        try:
            self.minio_client = Minio(
                endpoint=v.get_string("minio.endpoint"),
                access_key=v.get_string("minio.access_key_id"),
                secret_key=v.get_string("minio.access_key_secret"),
                secure=v.get_bool("minio.use_ssl"),
            )
        except Exception as e:
            self.logger.error(f"failed to initialize persistent storage client: {e}")
            raise FailedToInitializePersistentStorageError(error=e)

        self.minio_bucket_name = v.get_string("minio.bucket")
        if not self.minio_client.bucket_exists(self.minio_bucket_name):
            self.logger.error(f"bucket {self.minio_bucket_name} does not exist in persistent storage")
            self.minio_client = None
            raise MissingBucketError(self.minio_bucket_name)

        self.logger.debug(f"successfully initialized persistent storage with bucket {self.minio_bucket_name}!")

    def save(self, key: str, payload: BinaryIO, ttl_days: Optional[int] = None) -> str | None:
        try:
            result = None
            if ttl_days is not None:
                expiration_date = datetime.utcnow().replace(
                    hour=0,
                    minute=0,
                    second=0,
                    microsecond=0,
                ) + timedelta(days=ttl_days)
                result = self.minio_client.put_object(
                    self.minio_bucket_name,
                    key,
                    payload,
                    payload.getbuffer().nbytes,
                    retention=Retention(COMPLIANCE, expiration_date),
                )
            else:
                result = self.minio_client.put_object(
                    self.minio_bucket_name,
                    key,
                    payload,
                    payload.getbuffer().nbytes,
                )
            self.logger.info(f"file {key} successfully saved in persistent storage bucket {self.minio_bucket_name}")
            return result.version_id
        except Exception as e:
            error = FailedToSaveFileError(key, self.minio_bucket_name, e)
            self.logger.warning(f"{error}")
            raise error

    def get(self, key: str, version: Optional[str] = None) -> tuple[Optional[bytes], bool]:
        response = None
        try:
            exist = self._object_exist(key, version)
            if not exist:
                self.logger.error(
                    f"file {key} with version {version} not found in persistent storage bucket {self.minio_bucket_name}"
                )
                return None, False

            response = self.minio_client.get_object(self.minio_bucket_name, key, version_id=version)
            self.logger.info(
                f"file {key} successfully retrieved from persistent storage bucket {self.minio_bucket_name}"
            )
            return response.read(), True
        except Exception as e:
            error = FailedToGetFileError(key, version, self.minio_bucket_name, e)
            self.logger.error(f"{error}")
            raise error
        finally:
            if response:
                response.close()
                response.release_conn()

    def list(self) -> list[str]:
        try:
            objects = self.minio_client.list_objects(self.minio_bucket_name)
            self.logger.info(f"files successfully listed from persistent storage bucket {self.minio_bucket_name}")
            return [obj.object_name for obj in objects]
        except Exception as e:
            self.logger.error(FailedToListFilesError(self.minio_bucket_name, e))
            return []

    def list_versions(self, key: str) -> list[str]:
        try:
            objects = self.minio_client.list_objects(self.minio_bucket_name, prefix=key, include_version=True)
            self.logger.info(f"files successfully listed from persistent storage bucket {self.minio_bucket_name}")
            return [obj.object_name for obj in objects]
        except Exception as e:
            self.logger.error(f"failed to list files from persistent storage bucket {self.minio_bucket_name}: {e}")
            return []

    def delete(self, key: str, version: Optional[str] = None) -> bool:
        try:
            exist = self._object_exist(key, version)
            if not exist:
                self.logger.error(
                    f"file {key} with version {version} does not found in persistent storage bucket {self.minio_bucket_name}"
                )
                return False

            self.minio_client.remove_object(self.minio_bucket_name, key, version_id=version)
            self.logger.info(f"file {key} successfully deleted from persistent storage bucket {self.minio_bucket_name}")
            return True
        except Exception as e:
            error = FailedToDeleteFileError(key, version, self.minio_bucket_name, e)
            self.logger.error(f"{error}")
            raise error

    def _object_exist(self, key: str, version: str) -> bool:
        try:
            self.minio_client.stat_object(self.minio_bucket_name, key, version_id=version)
            return True
        except Exception as error:
            if "code: NoSuchKey" in str(error):
                return False
            else:
                raise error
