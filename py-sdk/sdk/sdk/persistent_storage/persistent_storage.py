from __future__ import annotations

import re
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from datetime import datetime, timedelta
from typing import BinaryIO, Optional

import loguru
from loguru import logger
from minio import Minio
from minio.commonconfig import ENABLED, Filter
from minio.credentials import ClientGrantsProvider
from minio.lifecycleconfig import Expiration, LifecycleConfig, Rule
from vyper import v

from sdk.auth.authentication import Authentication
from sdk.metadata.metadata import Metadata
from sdk.persistent_storage.exceptions import (
    FailedToDeleteFileError,
    FailedToGetFileError,
    FailedToInitializePersistentStorageError,
    FailedToListFilesError,
    FailedToSaveFileError,
    MissingBucketError,
)


@dataclass
class ObjectInfo:
    key: str = field(init=True)
    version: str = field(init=True)
    metadata: dict[str, str] = field(init=True)
    expires: Optional[datetime] = field(init=True)


@dataclass
class Object(ObjectInfo):
    data: bytes = field(init=True)


@dataclass
class PersistentStorageABC(ABC):
    @abstractmethod
    def save(self, key: str, payload: bytes, ttl_days: Optional[int]) -> Optional[ObjectInfo]:
        pass

    @abstractmethod
    def get(self, key: str, version: Optional[str]) -> Optional[Object]:
        pass

    @abstractmethod
    def list(self) -> list[ObjectInfo]:
        pass

    @abstractmethod
    def list_versions(self, key: str) -> list[ObjectInfo]:
        pass

    @abstractmethod
    def delete(self, key: str, version: Optional[str]) -> bool:
        pass


@dataclass
class PersistentStorage(PersistentStorageABC):
    logger: loguru.Logger = logger.bind(context="[PERSISTENT STORAGE]")
    minio_client: Minio = field(init=False)
    minio_bucket_name: str = field(init=False)
    metadata: Metadata = field(init=False, default_factory=Metadata)

    def __post_init__(self) -> None:
        try:
            print("Initializing Persistent Storage")
            auth = Authentication()

            print("Initializing Minio Creds")
            creds = ClientGrantsProvider(
                jwt_provider_func=lambda: auth.get_token(),
                sts_endpoint=f"{'https://' if v.get_bool('minio.ssl') else 'http://'}{v.get_string('minio.endpoint')}",
            )

            print("Initializing Minio Client")
            self.minio_client = Minio(
                endpoint=v.get_string("minio.endpoint"),
                credentials=creds,
                secure=v.get_bool("minio.ssl"),
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

    def save(self, key: str, payload: BinaryIO, ttl_days: Optional[int] = None) -> Optional[ObjectInfo]:
        try:
            if ttl_days is not None:
                rule = Rule(
                    rule_id=f"ttl-{key}",
                    status=ENABLED,
                    rule_filter=Filter(prefix=key),
                    expiration=Expiration(days=ttl_days),
                )
                conf = self.minio_client.get_bucket_lifecycle(self.minio_bucket_name)

                if conf is None:
                    conf = LifecycleConfig([rule])
                elif rule.rule_id not in {r.rule_id for r in conf.rules}:
                    conf.rules.append(rule)

                # The lifecycle should always be added before saving the object to the bucket
                self.minio_client.set_bucket_lifecycle(self.minio_bucket_name, conf)

            metadata = {
                "product": self.metadata.get_product(),
                "version": self.metadata.get_version(),
                "workflow": self.metadata.get_workflow(),
                "process": self.metadata.get_process(),
            }

            obj = self.minio_client.put_object(
                self.minio_bucket_name,
                key,
                payload,
                payload.getbuffer().nbytes,
                metadata=metadata,
            )
            self.logger.info(f"file {key} successfully saved in persistent storage bucket {self.minio_bucket_name}")

            expiry_date = self._get_expiry_date(obj.http_headers["x-amz-expiration"])

            return ObjectInfo(key=obj.object_name, version=obj.version_id, metadata=metadata, expires=expiry_date)
        except Exception as e:
            error = FailedToSaveFileError(key, self.minio_bucket_name, e)
            self.logger.warning(f"{error}")
            raise error

    def get(self, key: str, version: Optional[str] = None) -> Optional[Object]:
        response = None
        try:
            exist = self._object_exist(key, version)
            if not exist:
                self.logger.error(
                    f"file {key} with version {version} not found in persistent storage bucket {self.minio_bucket_name}"
                )
                return None

            response = self.minio_client.get_object(self.minio_bucket_name, key, version_id=version)
            self.logger.info(
                f"file {key} successfully retrieved from persistent storage bucket {self.minio_bucket_name}"
            )

            expiry_date = self._get_expiry_date(response.headers.get("x-amz-expiration"))

            return Object(key=key, version=response.headers.get("x-amz-version-id"), data=response.read(), expires=expiry_date, metadata=self._process_raw_metadata(response.headers))
        except Exception as e:
            error = FailedToGetFileError(key, version, self.minio_bucket_name, e)
            self.logger.error(f"{error}")
            raise error
        finally:
            if response:
                response.close()
                response.release_conn()

    def list(self) -> list[ObjectInfo]:
        try:
            objects = self.minio_client.list_objects(
                self.minio_bucket_name,
                recursive=True,
                include_user_meta=True,
                )
            self.logger.info(f"files successfully listed from persistent storage bucket {self.minio_bucket_name}")

            object_info_list = []

            # Get stats for each object that is not a directory
            for obj in objects:
                if not obj.is_dir:
                    stats = self.minio_client.stat_object(
                        self.minio_bucket_name, obj.object_name, version_id=obj.version_id
                    )

                    expiry_date = self._get_expiry_date(stats.metadata.get("x-amz-expiration"))

                    object_info_list.append(ObjectInfo(key=obj.object_name, version=stats.version_id, expires=expiry_date, metadata=self._process_raw_metadata(stats.metadata)))
            return object_info_list
        except Exception as e:
            self.logger.error(FailedToListFilesError(self.minio_bucket_name, e))
            return []

    def list_versions(self, key: str) -> list[ObjectInfo]:
        try:
            objects = self.minio_client.list_objects(
                self.minio_bucket_name,
                prefix=key,
                include_version=True,
                recursive=True,
                include_user_meta=True,
            )
            self.logger.info(f"files successfully listed from persistent storage bucket {self.minio_bucket_name}")

            object_info_list = []

            # Get stats for each object
            for obj in objects:
                if not obj.is_dir:
                    stats = self.minio_client.stat_object(
                        self.minio_bucket_name, obj.object_name, version_id=obj.version_id
                    )
                    expiry_date = self._get_expiry_date(stats.metadata.get("x-amz-expiration"))

                    object_info_list.append(ObjectInfo(key=obj.object_name, version=stats.version_id, expires=expiry_date, metadata=self._process_raw_metadata(stats.metadata)))

            return object_info_list
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
        # minio does not have a method to check if an object exists
        try:
            self.minio_client.stat_object(self.minio_bucket_name, key, version_id=version)
            return True
        except Exception as error:
            if "code: NoSuchKey" in str(error):
                return False
            else:
                raise error

    def _get_expiry_date(self, expiry_date_header: str) -> datetime | None:
        # Extract the date part using regular expression
        if not expiry_date_header:
            return None
        
        match = re.search(r'expiry-date="(.*?)"', expiry_date_header)
        if match:
            date_str = match.group(1)
            return datetime.strptime(date_str, "%a, %d %b %Y %H:%M:%S %Z")
        else:
            return None
        
    def _process_raw_metadata(self, raw_metadata: dict[str, str]) -> dict[str, str]:
        # minio replaces our metadata keys with x-amz-meta-<key>
        return {k.replace("x-amz-meta-", ""): v for k, v in raw_metadata.items() if k.startswith("x-amz-meta-")}
