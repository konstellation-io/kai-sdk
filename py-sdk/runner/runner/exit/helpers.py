from __future__ import annotations

from typing import TYPE_CHECKING, Optional

from google.protobuf.any_pb2 import Any

from runner.common.common import Finalizer, Handler, Initializer, initialize_process_configuration

if TYPE_CHECKING:
    from runner.exit.exit_runner import Preprocessor, Postprocessor

import inspect

from sdk.kai_sdk import KaiSDK


def compose_initializer(initializer: Optional[Initializer] = None) -> Initializer:
    async def initializer_func(sdk: KaiSDK) -> None:
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[INITIALIZER]")
        logger.info("initializing ExitRunner...")
        await sdk.initialize()
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


def compose_preprocessor(preprocessor: Preprocessor) -> Preprocessor:
    async def preprocessor_func(sdk: KaiSDK, response: Any) -> None:
        assert sdk.logger is not None
        
        product_id = sdk.metadata.get_product()
        version_id = sdk.metadata.get_version()
        workflow_id = sdk.metadata.get_workflow()
        process_id = sdk.metadata.get_process()
        request_id = sdk.get_request_id()
        metadata_info = f"{product_id=} {version_id=} {workflow_id=} {process_id=} {request_id=}"
        
        logger = sdk.logger.bind(context="[PREPROCESSOR]", metadata_info=metadata_info)
        logger.info("preprocessing ExitRunner...")

        logger.info("executing user preprocessor...")
        if inspect.iscoroutinefunction(preprocessor):
            await preprocessor(sdk, response)
        else:
            preprocessor(sdk, response)

    return preprocessor_func


def compose_handler(handler: Handler) -> Handler:
    async def handler_func(sdk: KaiSDK, response: Any) -> None:
        assert sdk.logger is not None

        product_id = sdk.metadata.get_product()
        version_id = sdk.metadata.get_version()
        workflow_id = sdk.metadata.get_workflow()
        process_id = sdk.metadata.get_process()
        request_id = sdk.get_request_id()
        metadata_info = f"{product_id=} {version_id=} {workflow_id=} {process_id=} {request_id=}"

        logger = sdk.logger.bind(context="[HANDLER]", metadata_info=metadata_info)
        logger.info("handling ExitRunner...")

        logger.info("executing user handler...")
        if inspect.iscoroutinefunction(handler):
            await handler(sdk, response)
        else:
            handler(sdk, response)

    return handler_func


def compose_postprocessor(postprocessor: Postprocessor) -> Postprocessor:
    async def postprocessor_func(sdk: KaiSDK, response: Any) -> None:
        assert sdk.logger is not None
        
        product_id = sdk.metadata.get_product()
        version_id = sdk.metadata.get_version()
        workflow_id = sdk.metadata.get_workflow()
        process_id = sdk.metadata.get_process()
        request_id = sdk.get_request_id()
        metadata_info = f"{product_id=} {version_id=} {workflow_id=} {process_id=} {request_id=}"
        
        logger = sdk.logger.bind(context="[POSTPROCESSOR]", metadata_info=metadata_info)
        logger.info("postprocessing ExitRunner...")

        logger.info("executing user postprocessor...")
        if inspect.iscoroutinefunction(postprocessor):
            await postprocessor(sdk, response)
        else:
            postprocessor(sdk, response)

    return postprocessor_func


def compose_finalizer(finalizer: Optional[Finalizer] = None) -> Finalizer:
    async def finalizer_func(sdk: KaiSDK) -> None:
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[FINALIZER]")
        logger.info("finalizing ExitRunner...")

        if finalizer is not None:
            logger.info("executing user finalizer...")
            if inspect.iscoroutinefunction(finalizer):
                await finalizer(sdk)
            else:
                finalizer(sdk)
            logger.info("user finalizer executed")

        logger.info("ExitRunner finalized")

    return finalizer_func
