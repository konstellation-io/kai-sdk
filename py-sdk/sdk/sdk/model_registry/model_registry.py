from __future__ import annotations

import os
import re
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import BinaryIO, Optional

import loguru
from loguru import logger
from minio import Minio
from minio.credentials import ClientGrantsProvider
from semver import Version
from vyper import v

from sdk.auth.authentication import Authentication
from sdk.metadata.metadata import Metadata
from sdk.model_registry.exceptions import (
    EmptyModelError,
    EmptyNameError,
    FailedToDeleteModelError,
    FailedToGetModelError,
    FailedToInitializeModelRegistryError,
    FailedToListModelsError,
    FailedToSaveModelError,
    InvalidVersionError,
    MissingBucketError,
)


@dataclass
class ModelInfo:
    name: str = field(init=True)
    version: str = field(init=True)
    description: str = field(init=True)
    format: str = field(init=True)


@dataclass
class Model(ModelInfo):
    model: BinaryIO = field(init=True)


@dataclass
class ModelRegistryABC(ABC):
    @abstractmethod
    def register_model(self, model: BinaryIO, name: str, version: str, description: str, model_format: str) -> None:
        pass

    @abstractmethod
    def get_model(self, name: str, version: Optional[str]) -> Optional[Model]:
        pass

    @abstractmethod
    def list_models(self) -> list[ModelInfo]:
        pass

    @abstractmethod
    def list_model_versions(self, name: str) -> list[ModelInfo]:
        pass

    @abstractmethod
    def delete_model(self, key: str) -> bool:
        pass


@dataclass
class ModelRegistry(ModelRegistryABC):
    logger: loguru.Logger = field(init=False)
    minio_client: Minio = field(init=False)
    minio_bucket_name: str = field(init=False)
    model_folder_name: str = field(init=False)

    def __post_init__(self) -> None:
        origin = logger._core.extra["origin"]
        self.logger = logger.bind(context=f"{origin}.[MODEL REGISTRY]")
        try:
            self.logger.info("initializing model registry")
            auth = Authentication()

            creds = ClientGrantsProvider(
                jwt_provider_func=lambda: auth.get_token(),
                sts_endpoint=f"{'https://' if v.get_bool('minio.ssl') else 'http://'}{v.get_string('minio.endpoint')}",
            )

            self.minio_client = Minio(
                endpoint=v.get_string("minio.endpoint"),
                credentials=creds,
                secure=v.get_bool("minio.ssl"),
            )
        except Exception as e:
            self.logger.error(f"failed to initialize model registry client: {e}")
            raise FailedToInitializeModelRegistryError(error=e)

        self.minio_bucket_name = v.get_string("minio.bucket")
        self.model_folder_name = os.path.join(
            v.get_string("minio.internal_folder"), v.get_string("model_registry.folder_name")
        )
        if not self.minio_client.bucket_exists(self.minio_bucket_name):
            self.logger.error(f"bucket {self.minio_bucket_name} does not exist in model registry")
            self.minio_client = None
            raise MissingBucketError(self.minio_bucket_name)

        self.logger.debug(f"successfully initialized model registry with bucket {self.minio_bucket_name}!")

    def register_model(self, model: BinaryIO, name: str, version: str, description: str, model_format: str) -> None:
        if name is None:
            raise EmptyNameError()

        if not Version.is_valid(version):
            raise InvalidVersionError()

        if not model:
            raise EmptyModelError()

        try:
            metadata = {
                "product": Metadata.get_product(),
                "version": Metadata.get_version(),
                "workflow": Metadata.get_workflow(),
                "process": Metadata.get_process(),
                "Model_version": version,
                "Model_description": description,
                "Model_format": model_format,
            }

            self.minio_client.put_object(
                bucket_name=self.minio_bucket_name,
                object_name=self._get_model_path(name),
                data=model,
                length=model.getbuffer().nbytes,
                metadata=metadata,
            )
            self.logger.info(f"model {name} successfully saved stored to the model registry with version {version}")
        except Exception as e:
            error = FailedToSaveModelError(name, e)
            self.logger.warning(f"{error}")
            raise error

    def get_model(self, name: str, version: Optional[str] = None) -> Optional[Model]:
        if name is None:
            raise EmptyNameError()

        if version and not Version.is_valid(version):
            raise InvalidVersionError()

        response = None
        try:
            exist = self._object_exist(self._get_model_path(name), version)
            if not exist:
                self.logger.error(f"model {name} with version {version} not found in model registry")
                return None

            response = self.minio_client.get_object(
                self.minio_bucket_name, self._get_model_path(name), version_id=version
            )
            self.logger.info(f"model {name} successfully retrieved from model registry")

            return Model(
                name=name,
                version=response.headers.get("x-amz-Model_version"),
                description=response.headers.get("x-amz-Model_description"),
                format=response.headers.get("x-amz-Model_format"),
                model=response.read(),
            )
        except Exception as e:
            error = FailedToGetModelError(name, version, e)
            self.logger.error(f"{error}")
            raise error
        finally:
            if response:
                response.close()
                response.release_conn()

    def list_models(self) -> list[ModelInfo]:
        try:
            objects = self.minio_client.list_objects(
                self.minio_bucket_name,
                recursive=True,
                include_user_meta=True,
                prefix=self._get_model_path(""),
            )
            self.logger.info("models successfully listed from model registry")

            model_info_list = []

            # Get stats for each object that is not a directory
            for obj in objects:
                if not obj.is_dir:
                    stats = self.minio_client.stat_object(
                        self.minio_bucket_name, obj.object_name, version_id=obj.version_id
                    )

                    model_info_list.append(
                        ModelInfo(
                            name=obj.object_name,
                            version=stats.metadata.get("x-amz-Model_version"),
                            format=stats.metadata.get("x-amz-Model_format"),
                            description=stats.metadata.get("x-amz-Model_description"),
                        )
                    )
            return model_info_list
        except Exception as e:
            self.logger.error(FailedToListModelsError(e))
            return []

    def list_model_versions(self, name: str) -> list[ModelInfo]:
        try:
            objects = self.minio_client.list_objects(
                self.minio_bucket_name,
                prefix=self._get_model_path(name),
                include_version=True,
                recursive=True,
                include_user_meta=True,
            )
            self.logger.info(f"models versions successfully listed from model registry for model {name}")
            model_version_info_list = []

            # Get stats for each object
            for obj in objects:
                if not obj.is_dir:
                    stats = self.minio_client.stat_object(
                        self.minio_bucket_name, obj.object_name, version_id=obj.version_id
                    )

                    model_version_info_list.append(
                        ModelInfo(
                            name=obj.object_name,
                            version=stats.metadata.get("x-amz-Model_version"),
                            format=stats.metadata.get("x-amz-Model_format"),
                            description=stats.metadata.get("x-amz-Model_description"),
                        )
                    )

            return model_version_info_list
        except Exception as e:
            self.logger.error(FailedToListModelsError(e))
            return []

    def delete_model(self, name: str) -> bool:
        try:
            exist = self._object_exist(name)
            if not exist:
                self.logger.error(f"model {name} does not found in model registry")
                return False

            self.minio_client.remove_object(self.minio_bucket_name, self._get_model_name(name))
            self.logger.info(f"model {name} successfully deleted from the model registry")
            return True
        except Exception as e:
            error = FailedToDeleteModelError(name, e)
            self.logger.error(f"{error}")
            raise error

    def _object_exist(self, key: str, version: Optional[str] = None) -> bool:
        # minio does not have a method to check if an object exists
        try:
            self.minio_client.stat_object(self.minio_bucket_name, key, version_id=version)
            return True
        except Exception as error:
            if "code: NoSuchKey" in str(error):
                return False
            else:
                raise error

    def _process_raw_metadata(self, raw_metadata: dict[str, str]) -> dict[str, str]:
        # minio replaces our metadata keys with x-amz-meta-<key>
        return {k.replace("x-amz-meta-", ""): v for k, v in raw_metadata.items() if k.startswith("x-amz-meta-")}

    def _get_model_path(self, model_name: str) -> str:
        return os.path.join(self.model_folder_name, model_name)

    def _get_model_name(self, full_path: str) -> str:
        _, file_name = os.path.split(full_path)
        return file_name
