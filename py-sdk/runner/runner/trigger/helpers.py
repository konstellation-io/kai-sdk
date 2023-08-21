from __future__ import annotations

import inspect
import signal
from typing import TYPE_CHECKING, Optional

from google.protobuf.any_pb2 import Any
from nats.aio.client import Client as NatsClient


from runner.common.common import Finalizer, Initializer, initialize_process_configuration

if TYPE_CHECKING:
    from runner.trigger.trigger_runner import ResponseHandler, RunnerFunc, TriggerRunner

from sdk.kai_sdk import KaiSDK
import threading


def compose_initializer(initializer: Initializer) -> Initializer:
    async def initializer_func(sdk: KaiSDK):
        logger = sdk.logger.bind(context="[RUNNER]")
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
        logger = sdk.logger.bind(context="[RUNNER]")
        logger.info("executing TriggerRunner...")
        shutdown_event = threading.Event()

        if user_runner is not None:
            logger.info("executing user runner...")
            user_runner(trigger_runner, sdk)
            logger.info("user runner executed")

        def shutdown_handler(sig, frame):
            shutdown_event.set()
            logger.info("shutting down runner...")
            logger.info("closing opened channels...")
            for request_id, channel in runner.response_channels.items():
                channel.close()
                logger.info(f"channel closed for request id {request_id}")

            logger.info("runnerFunc shutdown")

        user_runner_thread = threading.Thread(target=user_runner_thread)
        user_runner_thread.start()     

        signal.signal(signal.SIGINT, shutdown_handler)
        signal.signal(signal.SIGTERM, shutdown_handler)

        shutdown_event.wait()

        # TODO wait group done? waitgroup.set()?

    return runner_func


def get_response_handler(handlers: dict[str, NatsClient]) -> ResponseHandler:
    def response_handler_func(sdk: KaiSDK, response: Any):
        logger = sdk.logger.bind(context="[RESPONSE HANDLER]")
        request_id = sdk.get_request_id()
        logger.info(f"message received with request id {request_id}")

        response_handler = handlers.pop(request_id, None)
        if response_handler:
            response_handler.publish(response)
            return
        
        logger.debug(f"no response handler found for request id {request_id}")

    return response_handler_func


def compose_finalizer(user_finalizer: Finalizer) -> Finalizer:
    def finalizer_func(sdk: KaiSDK):
        logger = sdk.logger.bind(context="[FINALIZER]")
        logger.info("finalizing TriggerRunner...")

        if user_finalizer is not None:
            logger.info("executing user finalizer...")
            user_finalizer(sdk)
            logger.info("user finalizer executed")

        logger.info("TriggerRunner finalized")

    return finalizer_func
