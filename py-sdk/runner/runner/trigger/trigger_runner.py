from __future__ import annotations

import asyncio
import concurrent.futures
import functools
import signal
import sys
import threading
from concurrent.futures import ThreadPoolExecutor
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
    tasks: list[threading.Thread] = field(init=False, default_factory=list)

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
        signal: Optional[signal.Signals] = None,
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

        await self.finalizer(self.sdk)
        self.logger.info("successfully shutdown trigger runner")

        [task.cancel() for task in self.tasks]

        self.logger.info(f"cancelling {len(self.tasks)} outstanding tasks")
        await asyncio.gather(*self.tasks, return_exceptions=True)

        self.logger.info("shutting down executor")
        executor.shutdown(wait=False)

        self.logger.info(f"releasing {len(executor._threads)} threads from executor")
        for thread in executor._threads:
            try:
                thread._tstate_lock.release()
            except Exception as e:
                self.logger.error(f"error releasing thread lock: {e}")

        if not self.nc.is_closed:
            self.logger.info("closing nats connection")
            await self.nc.close()

        loop.stop()

    def _exception_handler(self, loop, executor, context) -> None:
        msg = context.get("exception", context["message"])
        self.logger.error(f"caught exception: {msg}")
        asyncio.create_task(self._shutdown_handler(loop, executor))

    def subscriber_wrapper(self) -> None:
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        loop.run_until_complete(self.subscriber.start())

    def runner_wrapper(self) -> None:
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        loop.run_until_complete(self.runner(self, self.sdk))

    async def run(self) -> None:
        if getattr(self, "runner", None) is None:
            raise UndefinedRunnerFunctionError

        if not self.initializer:
            self.initializer = compose_initializer()

        self.response_handler = get_response_handler(self.response_channels)

        if not self.finalizer:
            self.finalizer = compose_finalizer()

        await self.initializer(self.sdk)

        loop = asyncio.get_event_loop()
        executor = ThreadPoolExecutor(max_workers=2)
        signals = (signal.SIGINT, signal.SIGTERM)
        for s in signals:
            loop.add_signal_handler(
                s,
                lambda s=s: asyncio.create_task(self._shutdown_handler(loop, executor, signal=s)),
            )
        exception_func_handler = functools.partial(self._exception_handler, executor)
        loop.set_exception_handler(exception_func_handler)

        try:
            future_run = loop.run_in_executor(executor, self.runner_wrapper)
            future_sub = loop.run_in_executor(executor, self.subscriber_wrapper)

            self.tasks = [future_run, future_sub]
            await asyncio.gather(*self.tasks, return_exceptions=True)
        finally:
            self.logger.info("trigger runner stopped")


RunnerFunc = Callable[[TriggerRunner, KaiSDK], Awaitable[None]]
