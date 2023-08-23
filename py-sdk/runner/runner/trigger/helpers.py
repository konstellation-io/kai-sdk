from __future__ import annotations

import inspect
from queue import Queue
from signal import SIGINT, SIGTERM, signal
from typing import TYPE_CHECKING

from google.protobuf.any_pb2 import Any

from runner.common.common import Finalizer, Initializer, initialize_process_configuration

if TYPE_CHECKING:
    from runner.trigger.trigger_runner import ResponseHandler, RunnerFunc, TriggerRunner

from sdk.kai_sdk import KaiSDK


def compose_initializer(initializer: Initializer) -> Initializer:
    async def initializer_func(sdk: KaiSDK):
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[INITIALIZER]")
        logger.info("initializing TriggerRunner...")
        await initialize_process_configuration(sdk)

        if initializer is not None:
            logger.info("executing user initializer...")
            if inspect.iscoroutinefunction(initializer):
                await initializer(sdk)
            else:
                initializer(sdk)
            logger.info("user initializer executed")

        logger.info("TriggerRunner initialized")

    return initializer_func


def compose_runner(trigger_runner: TriggerRunner, user_runner: RunnerFunc) -> RunnerFunc:
    def runner_func(runner: TriggerRunner, sdk: KaiSDK):
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[RUNNER]")
        logger.info("executing TriggerRunner...")

        if user_runner is not None:
            logger.info("executing user runner...")
            user_runner(trigger_runner, sdk)
            logger.info("user runner executed")

        async def shutdown_handler(sig, frame):
            logger.info("shutting down runner...")
            logger.info("closing opened channels...")
            for request_id, channel in runner.response_channels.items():
                channel.put(None)
                logger.info(f"channel closed for request id {request_id}")

            trigger_runner.runner_thread_shutdown_event.set()

        signal(SIGINT, shutdown_handler)
        signal(SIGTERM, shutdown_handler)

        trigger_runner.runner_thread_shutdown_event.wait()
        logger.info("runnerFunc shutdown")

    return runner_func


def get_response_handler(handlers: dict[str, Queue]) -> ResponseHandler:
    def response_handler_func(sdk: KaiSDK, response: Any):
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[RESPONSE HANDLER]")
        request_id = sdk.get_request_id()
        logger.info(f"message received with request id {request_id}")

        handler = handlers.pop(request_id, None)
        if handler:
            handler.put(response)
            return

        logger.debug(f"no response handler found for request id {request_id}")

    return response_handler_func


def compose_finalizer(user_finalizer: Finalizer) -> Finalizer:
    def finalizer_func(sdk: KaiSDK):
        assert sdk.logger is not None
        logger = sdk.logger.bind(context="[FINALIZER]")
        logger.info("finalizing TriggerRunner...")

        if user_finalizer is not None:
            logger.info("executing user finalizer...")
            user_finalizer(sdk)
            logger.info("user finalizer executed")

        logger.info("TriggerRunner finalized")

    return finalizer_func
