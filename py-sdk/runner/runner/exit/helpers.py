from __future__ import annotations

import inspect
from typing import TYPE_CHECKING, Optional

from google.protobuf.any_pb2 import Any

from runner.common.common import Finalizer, Initializer, initialize_process_configuration

if TYPE_CHECKING:
    from runner.exit.exit_runner import Preprocessor, Handler, Postprocessor


from sdk.kai_sdk import KaiSDK


def compose_initializer(initializer: Optional[Initializer] = None) -> Initializer:
    async def initializer_func(sdk: KaiSDK) -> None:
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[INITIALIZER]")
        logger.info("initializing ExitRunner...")
        await initialize_process_configuration(sdk)

        if initializer is not None:
            logger.info("executing user initializer...")
            await initializer(sdk)
            logger.info("user initializer executed")

        logger.info("ExitRunner initialized")

    return initializer_func


def compose_preprocessor(preprocessor: Preprocessor) -> Preprocessor:
    def preprocessor_func(sdk: KaiSDK, response: Any) -> None:
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[PREPROCESSOR]")
        logger.info("preprocessing ExitRunner...")

        logger.info("executing user preprocessor...")
        preprocessor(sdk, response)

    return preprocessor_func


def compose_handler(handler: Handler) -> Handler:
    def handler_func(sdk: KaiSDK, response: Any) -> None:
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[HANDLER]")
        logger.info("handling ExitRunner...")

        logger.info("executing user handler...")
        handler(sdk, response)

    return handler_func


def compose_postprocessor(postprocessor: Postprocessor) -> Postprocessor:
    def postprocessor_func(sdk: KaiSDK, response: Any) -> None:
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[POSTPROCESSOR]")
        logger.info("postprocessing ExitRunner...")

        logger.info("executing user postprocessor...")
        postprocessor(sdk, response)

    return postprocessor_func


def compose_finalizer(finalizer: Optional[Finalizer] = None) -> Finalizer:
    def finalizer_func(sdk: KaiSDK) -> None:
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[FINALIZER]")
        logger.info("finalizing ExitRunner...")

        if finalizer is not None:
            logger.info("executing user finalizer...")
            finalizer(sdk)
            logger.info("user finalizer executed")

        logger.info("ExitRunner finalized")

    return finalizer_func
