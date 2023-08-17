import signal
from typing import Optional

from google.protobuf.any_pb2 import Any

from runner.common.common import Finalizer, Initializer, initialize_process_configuration
from runner.trigger.trigger_runner import ResponseHandler, RunnerFunc, TriggerRunner
from sdk.kai_sdk import KaiSDK


def compose_initializer(initializer: Initializer) -> Initializer:
    async def initializer_func(sdk: KaiSDK):
        logger = sdk.logger.bind(context="[RUNNER]")
        logger.info("initializing TriggerRunner...")
        await initialize_process_configuration(sdk)

        if initializer is not None:
            logger.info("executing user initializer...")
            initializer(sdk)
            logger.info("user initializer executed")

        logger.info("TriggerRunner initialized")

    return initializer_func


def compose_runner(trigger_runner: TriggerRunner, user_runner: RunnerFunc) -> RunnerFunc:
    def runner_func(runner: TriggerRunner, sdk: KaiSDK):
        logger = sdk.logger.bind(context="[RUNNER]")
        logger.info("executing TriggerRunner...")

        if user_runner is not None:
            logger.info("executing user runner...")
            user_runner(trigger_runner, sdk)
            logger.info("[RUNNER] User runner executed")

        def shutdown_handler(signum, frame):
            runner.shutdown_event.set()

        signal.signal(signal.SIGINT, shutdown_handler)
        signal.signal(signal.SIGTERM, shutdown_handler)

        logger.info("shutting down runner...")
        logger.info("closing opened channels...")
        # for key, value in runner.response_channels.items():
        #     value.close()
        #     logger.info("channel closed for requestID {key}")

        runner.shutdown_event.wait()
        logger.info("runnerFunc shutdown")

    return runner_func


def get_response_handler(handlers) -> ResponseHandler:  # TODO define type
    def response_handler_func(sdk: KaiSDK, response: Any):
        logger = sdk.logger.bind(context="[RESPONSE HANDLER]")
        request_id = sdk.get_request_id()
        logger.info(f"message received with requestID {request_id}")

        # response_handler = handlers.pop(sdk.request_id, None)
        # if response_handler is not None:
        #     response_handler.put(response)
        # else:
        #     logger.info(f"no handler found for the message with requestID {request_id}")

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
