from __future__ import annotations

import asyncio
import concurrent.futures
import functools
import signal
import sys
from dataclasses import dataclass, field
from queue import Queue
from typing import Any, Awaitable, Callable, Optional

import loguru
from google.protobuf import any_pb2
from loguru import logger
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext

from runner.common.common import Finalizer, Initializer
from runner.trigger.exceptions import UndefinedRunnerFunctionError
from runner.trigger.helpers import compose_finalizer, compose_initializer, compose_runner, get_response_handler
from runner.trigger.subscriber import TriggerSubscriber
from sdk.kai_sdk import KaiSDK

ResponseHandler = Callable[[KaiSDK, any_pb2.Any], None]


@dataclass
class TriggerRunner:
    sdk: KaiSDK = field(init=False)
    nc: NatsClient
    js: JetStreamContext
    logger: loguru.Logger = logger.bind(context="[TRIGGER]")
    response_handler: ResponseHandler = field(init=False, default=None)
    response_channels: dict[str, Queue] = field(init=False, default_factory=dict)
    initializer: Optional[Initializer] = None
    runner: RunnerFunc = field(init=False)
    subscriber: TriggerSubscriber = field(init=False)
    finalizer: Optional[Finalizer] = None

    def __post_init__(self) -> None:
        self.sdk = KaiSDK(nc=self.nc, js=self.js, logger=self.logger)
        self.subscriber = TriggerSubscriber(self)

    def with_initializer(self, initializer: Initializer) -> TriggerRunner:
        self.initializer = compose_initializer(initializer)
        return self

    def with_runner(self, runner: RunnerFunc) -> TriggerRunner:
        self.runner = compose_runner(runner)
        return self

    def with_finalizer(self, finalizer: Finalizer) -> TriggerRunner:
        self.finalizer = compose_finalizer(finalizer)
        return self

    def get_response_channel(self, request_id: str) -> Any:
        if request_id not in self.response_channels:
            self.response_channels[request_id] = Queue(maxsize=1)
        return self.response_channels[request_id].get()

    async def _shutdown_handler(
        self,
        loop: asyncio.AbstractEventLoop,
        executor: concurrent.futures.ThreadPoolExecutor,
        signal: int = None,
    ) -> None:
        if signal:
            self.logger.info(f"received exit signal {signal.name}...")
        self.logger.info("shutting down runner...")
        self.logger.info("closing opened channels...")
        for request_id, channel in self.response_channels.items():
            channel.put(None)
            self.logger.info(f"channel closed for request id {request_id}")

        self.logger.info("shutting down subscriber")
        for sub in self.subscriber.subscriptions:
            self.logger.info(f"unsubscribing from subject {sub.subject}")

            try:
                await sub.unsubscribe()
            except Exception as e:
                self.logger.error(f"error unsubscribing from the NATS subject {sub.subject}: {e}")
                sys.exit(1)

        self.finalizer(self.sdk)
        self.logger.info("successfully shutdown trigger runner")

        tasks = [t for t in asyncio.all_tasks() if t is not asyncio.current_task()]

        [task.cancel() for task in tasks]

        self.logger.info(f"cancelling {len(tasks)} outstanding tasks")
        await asyncio.gather(*tasks, return_exceptions=True)

        self.logger.info("shutting down executor")
        executor.shutdown(wait=False)

        self.logger.info(f"releasing {len(executor._threads)} threads from executor")
        for thread in executor._threads:
            try:
                thread._tstate_lock.release()
            except Exception as e:
                self.logger.error(f"error releasing thread lock: {e}")

        loop.stop()

    def _exception_handler(self, executor, loop, context) -> None:
        msg = context.get("exception", context["message"])
        self.logger.error(f"caught exception: {msg}")
        asyncio.create_task(self._shutdown_handler(loop, executor))

    def runner_sync(self, loop):
        asyncio.run_coroutine_threadsafe(self.runner(self, self.sdk), loop)

    async def runner_coro(self, executor: concurrent.futures.ThreadPoolExecutor) -> None:
        loop = asyncio.get_running_loop()
        await loop.run_in_executor(executor, self.runner_sync, loop)

    def subscriber_sync(self, loop):
        asyncio.run_coroutine_threadsafe(self.subscriber.start(), loop)

    async def subscriber_coro(self, executor: concurrent.futures.ThreadPoolExecutor) -> None:
        loop = asyncio.get_running_loop()
        await loop.run_in_executor(executor, self.subscriber_sync, loop)

    async def run(self) -> None:
        if getattr(self, "runner", None) is None:
            raise UndefinedRunnerFunctionError

        if not self.initializer:
            self.initializer = compose_initializer()

        self.response_handler = get_response_handler(self.response_channels)

        if not self.finalizer:
            self.finalizer = compose_finalizer()

        initializer_func = self.initializer(self.sdk)
        await initializer_func

        executor = concurrent.futures.ThreadPoolExecutor()
        loop = asyncio.get_event_loop()
        signals = (signal.SIGINT, signal.SIGTERM)
        for s in signals:
            loop.add_signal_handler(
                s,
                lambda s=s: asyncio.create_task(self._shutdown_handler(loop, executor, signal=s)),
            )
        handle_exec_func = functools.partial(self._exception_handler, executor)
        loop.set_exception_handler(handle_exec_func)

        try:
            loop.create_task(self.runner_coro(executor))
            loop.create_task(self.subscriber_coro(executor))
        finally:
            self.logger.info("runner started")


RunnerFunc = Callable[[TriggerRunner, KaiSDK], Awaitable[None]]
