from typing import Awaitable, Callable, Optional

from google.protobuf.any_pb2 import Any
from vyper import v

from sdk.kai_sdk import KaiSDK

Task = Callable[[KaiSDK], None]
Initializer = Callable[[KaiSDK], Awaitable[None] | None]
Finalizer = Task
Handler = Callable[[KaiSDK, Any], None]


async def initialize_process_configuration(sdk: KaiSDK):
    values = v.get("centralized_configuration.process.config")

    logger = sdk.logger.bind(context="[CONFIG INITIALIZER]")
    logger.info("initializing process configuration")

    for key, value in values.items():
        try:
            await sdk.centralized_config.set_config(key, value)
        except Exception as e:
            logger.error(f"error initializing process configuration with key {key}: {e}")

    logger.info("process configuration initialized")
