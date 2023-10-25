from __future__ import annotations

from abc import ABC, abstractmethod
from dataclasses import dataclass

import loguru
from loguru import logger
from minio import Minio
from sdk.persistent_storage.exceptions import FailedPersistentStorageInitializationError
from vyper import v


@dataclass
class PersistentStorageABC(ABC):
    @abstractmethod
    def save(self, key: str, payload: bytes, ttl: int) -> None:
        pass

    @abstractmethod
    def get(self, key: str, version: str) -> tuple[bytes, bool]:
        pass

    @abstractmethod
    def list(self) -> list[str]:
        pass

    @abstractmethod
    def list(self, key: str) -> list[str]:
        pass

    @abstractmethod
    def delete(self, key: str, version: str) -> bool:
        pass


@dataclass
class PersistentStorage(PersistentStorageABC):
    logger: loguru.Logger = logger.bind(context="[PERSISTENT STORAGE]")
    minio_client: Minio = None

    def __post_init__(self) -> None:
        try:
            self.minio_client = Minio(
                endpoint=v.get_string("minio_endpoint"),
                access_key=v.get_string("minio_access_key"),
                secret_key=v.get_string("minio_secret_key"),
                region=v.get_string("minio_region"),
                secure=v.get_bool("minio_secure"),
            )
        except Exception as e:
            self.logger.error(f"Failed to initialize minio client: {e}")
            raise FailedPersistentStorageInitializationError(error=e)

    def save(self, key: str, payload: bytes, ttl: int = 30) -> None:
        pass

    def get(self, key: str, version: str = "latest") -> tuple[bytes, bool]:
        pass

    def list(self) -> list[str]:
        pass

    def list(self, key: str) -> list[str]:
        pass

    def delete(self, key: str, version: str = "latest") -> bool:
        pass
