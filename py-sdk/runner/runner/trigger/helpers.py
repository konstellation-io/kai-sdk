from __future__ import annotations

from queue import Queue
from signal import SIGINT, SIGTERM
from threading import Event
from typing import TYPE_CHECKING, Optional

from google.protobuf.any_pb2 import Any

from runner.common.common import Finalizer, Initializer, initialize_process_configuration

if TYPE_CHECKING:
    from runner.trigger.trigger_runner import ResponseHandler, RunnerFunc, TriggerRunner

from sdk.kai_sdk import KaiSDK


def compose_initializer(initializer: Optional[Initializer] = None) -> Initializer:
    async def initializer_func(sdk: KaiSDK) -> None:
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[INITIALIZER]")
        logger.info("initializing TriggerRunner...")
        await sdk.initialize()
        await initialize_process_configuration(sdk)

        if initializer is not None:
            logger.info("executing user initializer...")
            await initializer(sdk)
            logger.info("user initializer executed")

        logger.info("TriggerRunner initialized")

    return initializer_func


def compose_runner(user_runner: RunnerFunc) -> RunnerFunc:
    async def runner_func(trigger_runner: TriggerRunner, sdk: KaiSDK) -> None:
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[RUNNER]")
        logger.info("executing TriggerRunner...")

        logger.info("executing user runner...")
        await user_runner(trigger_runner, sdk)
        logger.info("user runner executed")

        logger.info("runnerFunc shutdown")

    return runner_func


def get_response_handler(handlers: dict[str, Queue]) -> ResponseHandler:
    def response_handler_func(sdk: KaiSDK, response: Any) -> None:
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[RESPONSE HANDLER]")
        request_id = sdk.get_request_id()
        assert request_id is not None
        logger.info(f"message received with request id {request_id}")

        handler = handlers.pop(request_id, None)
        if handler:
            handler.put(response)
            return

        logger.debug(f"no response handler found for request id {request_id}")

    return response_handler_func


def compose_finalizer(user_finalizer: Optional[Finalizer] = None) -> Finalizer:
    def finalizer_func(sdk: KaiSDK) -> None:
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[FINALIZER]")
        logger.info("finalizing TriggerRunner...")

        if user_finalizer is not None:
            logger.info("executing user finalizer...")
            user_finalizer(sdk)
            logger.info("user finalizer executed")

        logger.info("TriggerRunner finalized")

    return finalizer_func
