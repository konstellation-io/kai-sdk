import signal
from typing import Optional

from google.protobuf.any_pb2 import Any

from runner.common.common import Finalizer, Handler, Initializer, Task, initialize_process_configuration
from sdk.kai_sdk import KaiSDK


def compose_initializer(initializer: Initializer) -> Initializer:
    def initializer_func(sdk: KaiSDK):
        sdk.Logger.info("[RUNNER] Initializing TriggerRunner...")
        initialize_process_configuration(sdk)

        if initializer:
            sdk.Logger.info("[RUNNER] Executing user initializer...")
            initializer(sdk)
            sdk.Logger.info("[RUNNER] User initializer executed")
        sdk.Logger.info("[RUNNER] TriggerRunner initialized")

    return initializer_func


def compose_runner(trigger_runner: Runner, user_runner: RunnerFunc) -> RunnerFunc:
    def runner_func(runner: Runner, sdk: KaiSDK):
        sdk.Logger.info("[RUNNER] Running TriggerRunner...")

        if user_runner:
            sdk.Logger.info("[RUNNER] Executing user runner...")
            user_runner(trigger_runner, sdk)
            sdk.Logger.info("[RUNNER] User runner executed")

        def shutdown_handler(signum, frame):
            runner.shutdown_event.set()

        signal.signal(signal.SIGINT, shutdown_handler)
        signal.signal(signal.SIGTERM, shutdown_handler)

        sdk.Logger.info("[RUNNER] Shutting down runner...")
        sdk.Logger.info("[RUNNER] Closing opened channels...")
        for key, value in runner.response_channels.items():
            value.close()
            sdk.Logger.info("[RUNNER] Channel closed for requestID", "RequestID", key)

        runner.shutdown_event.wait()
        sdk.Logger.info("[RUNNER] RunnerFunc shutdown")

    return runner_func


def get_response_handler(handlers) -> ResponseHandler:
    def response_handler(sdk: KaiSDK, response: Any) -> Optional[Exception]:
        sdk.Logger.info("[RESPONSE HANDLER] Message received", "RequestID", sdk.request_id)

        response_handler = handlers.pop(sdk.request_id, None)
        if response_handler:
            response_handler.put(response)
        else:
            sdk.Logger.info("[RESPONSE HANDLER] No handler found for the message", "RequestID", sdk.request_id)

    return response_handler


def compose_finalizer(user_finalizer: Finalizer) -> Finalizer:
    def finalizer_func(sdk: KaiSDK):
        sdk.Logger.info("[FINALIZER] Finalizing TriggerRunner...")

        if user_finalizer:
            sdk.Logger.info("[FINALIZER] Executing user finalizer...")
            user_finalizer(sdk)
            sdk.Logger.info("[FINALIZER] User finalizer executed")

        sdk.Logger.info("[FINALIZER] TriggerRunner finalized")

    return finalizer_func
