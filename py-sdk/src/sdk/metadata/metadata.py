from dataclasses import dataclass
from loguru import logger
from loguru._logger import Logger
from vyper import v

@dataclass
class Metadata:
    logger: Logger = logger.bind(context="[METADATA]")

    @staticmethod
    def get_product() -> str:
        return v.get("metadata.product_id")

    @staticmethod
    def get_workflow() -> str:
        return v.get("metadata.workflow_id")

    @staticmethod
    def get_process() -> str:
        return v.get("metadata.process_id")

    @staticmethod
    def get_version() -> str:
        return v.get("metadata.version_id")

    @staticmethod
    def get_object_store_name() -> str:
        return v.get("nats.object_store")

    @staticmethod
    def get_key_value_store_product_name() -> str:
        return v.get("centralized_configuration.product.bucket")

    @staticmethod
    def get_key_value_store_workflow_name() -> str:
        return v.get("centralized_configuration.workflow.bucket")

    @staticmethod
    def get_key_value_store_process_name() -> str:
        return v.get("centralized_configuration.process.bucket")
