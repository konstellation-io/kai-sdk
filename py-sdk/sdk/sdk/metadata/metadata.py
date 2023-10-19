from __future__ import annotations

from abc import ABC, abstractmethod
from dataclasses import dataclass

import loguru
from loguru import logger
from vyper import v


@dataclass
class MetadataABC(ABC):
    @staticmethod
    @abstractmethod
    def get_global() -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_product() -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_workflow() -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_process() -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_version() -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_ephemeral_storage_name() -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_key_value_global_name() -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_key_value_store_product_name() -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_key_value_store_workflow_name() -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_key_value_store_process_name() -> str:
        pass


@dataclass
class Metadata(MetadataABC):
    logger: loguru.Logger = logger.bind(context="[METADATA]")

    @staticmethod
    def get_global() -> str:
        return v.get_string("metadata.global_id")

    @staticmethod
    def get_product() -> str:
        return v.get_string("metadata.product_id")

    @staticmethod
    def get_workflow() -> str:
        return v.get_string("metadata.workflow_id")

    @staticmethod
    def get_process() -> str:
        return v.get_string("metadata.process_id")

    @staticmethod
    def get_version() -> str:
        return v.get_string("metadata.version_id")

    @staticmethod
    def get_ephemeral_storage_name() -> str:
        return v.get_string("nats.object_store")

    @staticmethod
    def get_key_value_global_name() -> str:
        return v.get_string("centralized_configuration.global.bucket")

    @staticmethod
    def get_key_value_store_product_name() -> str:
        return v.get_string("centralized_configuration.product.bucket")

    @staticmethod
    def get_key_value_store_workflow_name() -> str:
        return v.get_string("centralized_configuration.workflow.bucket")

    @staticmethod
    def get_key_value_store_process_name() -> str:
        return v.get_string("centralized_configuration.process.bucket")
