from typing import Callable, Optional

from google.protobuf.any_pb2 import Any
from vyper import v

from sdk.kai_sdk import KaiSDK

Task = Callable[[KaiSDK], None]
Initializer = Task
Finalizer = Task
Handler = Callable[[KaiSDK, Any], Optional[Exception]]


async def initialize_process_configuration(sdk: KaiSDK):
    values = v._get_key_value_config("centralized_configuration.process.config")

    logger = sdk.logger.bind("[CONFIG INITIALIZER]")
    logger.info("initializing process configuration")

    for key, value in values.items():
        try:
            await sdk.centralized_config.set_config(key, value)
        except Exception as e:
            logger.error(f"error initializing process configuration with key {key}: {e}")

    logger.info("process configuration initialized")
