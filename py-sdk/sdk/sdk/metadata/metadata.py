from abc import ABC, abstractmethod
from dataclasses import dataclass

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
    def get_object_store_name() -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_global_centralized_configuration_name() -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_product_centralized_configuration_name() -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_workflow_centralized_configuration_name() -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_process_centralized_configuration_name() -> str:
        pass


@dataclass
class Metadata(MetadataABC):
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
    def get_object_store_name() -> str:
        return v.get_string("nats.object_store")

    @staticmethod
    def get_global_centralized_configuration_name() -> str:
        return v.get_string("centralized_configuration.global.bucket")

    @staticmethod
    def get_product_centralized_configuration_name() -> str:
        return v.get_string("centralized_configuration.product.bucket")

    @staticmethod
    def get_workflow_centralized_configuration_name() -> str:
        return v.get_string("centralized_configuration.workflow.bucket")

    @staticmethod
    def get_process_centralized_configuration_name() -> str:
        return v.get_string("centralized_configuration.process.bucket")
