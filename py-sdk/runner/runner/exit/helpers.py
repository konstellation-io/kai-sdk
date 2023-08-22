from __future__ import annotations

import inspect
from typing import TYPE_CHECKING, Optional

import loguru
from google.protobuf.any_pb2 import Any

from runner.common.common import Finalizer, Initializer, initialize_process_configuration

if TYPE_CHECKING:
    from runner.exit.exit_runner import Preprocessor, Handler, Postprocessor

import loguru
from google.protobuf.any_pb2 import Any

from sdk.kai_sdk import KaiSDK


def compose_initializer(initializer: Initializer) -> Initializer:
    async def initializer_func(sdk: KaiSDK):
        assert isinstance(sdk.logger, loguru.Logger)
        logger = sdk.logger.bind(context="[INITIALIZER]")
        logger.info("initializing ExitRunner...")
        await initialize_process_configuration(sdk)

        if initializer is not None:
            logger.info("executing user initializer...")
            if inspect.iscoroutinefunction(initializer):
                await initializer(sdk)
            else:
                initializer(sdk)
            logger.info("user initializer executed")

        logger.info("ExitRunner initialized")

    return initializer_func


def compose_preprocessor(preprocessor: Preprocessor) -> Optional[Preprocessor]:
    def preprocessor_func(sdk: KaiSDK, response: Any):
        assert isinstance(sdk.logger, loguru.Logger)
        logger = sdk.logger.bind(context="[PREPROCESSOR]")
        logger.info("preprocessing ExitRunner...")

        if preprocessor is not None:
            logger.info("executing user preprocessor...")
            preprocessor(sdk, response)

        return None

    return preprocessor_func


def compose_handler(handler: Handler) -> Optional[Handler]:
    def handler_func(sdk: KaiSDK, response: Any):
        assert isinstance(sdk.logger, loguru.Logger)
        logger = sdk.logger.bind(context="[HANDLER]")
        logger.info("handling ExitRunner...")

        if handler is not None:
            logger.info("executing user handler...")
            handler(sdk, response)

        return None

    return handler_func


def compose_postprocessor(postprocessor: Postprocessor) -> Optional[Postprocessor]:
    def postprocessor_func(sdk: KaiSDK, response: Any):
        assert isinstance(sdk.logger, loguru.Logger)
        logger = sdk.logger.bind(context="[POSTPROCESSOR]")
        logger.info("postprocessing ExitRunner...")

        if postprocessor is not None:
            logger.info("executing user postprocessor...")
            postprocessor(sdk, response)

        return None

    return postprocessor_func


def compose_finalizer(finalizer: Finalizer) -> Finalizer:
    def finalizer_func(sdk: KaiSDK):
        assert isinstance(sdk.logger, loguru.Logger)
        logger = sdk.logger.bind(context="[FINALIZER]")
        logger.info("finalizing ExitRunner...")

        if finalizer is not None:
            logger.info("executing user finalizer...")
            finalizer(sdk)
            logger.info("user finalizer executed")

        logger.info("ExitRunner finalized")

    return finalizer_func
