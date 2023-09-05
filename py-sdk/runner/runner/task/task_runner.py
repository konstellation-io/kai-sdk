from __future__ import annotations

import asyncio
import signal
import sys
from dataclasses import dataclass, field
from typing import Any, Awaitable, Callable, Optional

import loguru
from loguru import logger
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext

from runner.common.common import Finalizer, Handler, Initializer
from runner.task.exceptions import UndefinedDefaultHandlerFunctionError
from runner.task.helpers import (
    compose_finalizer,
    compose_handler,
    compose_initializer,
    compose_postprocessor,
    compose_preprocessor,
)
from runner.task.subscriber import TaskSubscriber
from sdk.kai_sdk import KaiSDK

Preprocessor = Callable[[KaiSDK, Any], Awaitable[None] | None]
Postprocessor = Callable[[KaiSDK, Any], Awaitable[None] | None]


@dataclass
class TaskRunner:
    sdk: KaiSDK = field(init=False)
    nc: NatsClient
    js: JetStreamContext
    logger: loguru.Logger = logger.bind(context="[TASK]")
    response_handlers: dict[str, Handler] = field(default_factory=dict)
    initializer: Optional[Initializer] = None
    preprocessor: Optional[Preprocessor] = None
    postprocessor: Optional[Postprocessor] = None
    finalizer: Optional[Finalizer] = None

    def __post_init__(self) -> None:
        self.sdk = KaiSDK(nc=self.nc, js=self.js, logger=self.logger)
        self.subscriber = TaskSubscriber(self)

    def with_initializer(self, initializer: Initializer) -> TaskRunner:
        self.initializer = compose_initializer(initializer)
        return self

    def with_preprocessor(self, preprocessor: Preprocessor) -> TaskRunner:
        self.preprocessor = compose_preprocessor(preprocessor)
        return self

    def with_handler(self, handler: Handler) -> TaskRunner:
        self.response_handlers["default"] = compose_handler(handler)
        return self

    def with_custom_handler(self, subject: str, handler: Handler) -> TaskRunner:
        self.response_handlers[subject] = compose_handler(handler)
        return self

    def with_postprocessor(self, postprocessor: Postprocessor) -> TaskRunner:
        self.postprocessor = compose_postprocessor(postprocessor)
        return self

    def with_finalizer(self, finalizer: Finalizer) -> TaskRunner:
        self.finalizer = compose_finalizer(finalizer)
        return self

    def _exception_handler(self, loop, context) -> None:
        msg = context.get("exception", context["message"])
        self.logger.error(f"caught exception: {msg}")
        asyncio.create_task(self._shutdown_handler(loop))

    async def _shutdown_handler(self, loop: asyncio.AbstractEventLoop, signal: int = None) -> None:
        if signal:
            self.logger.info(f"received exit signal {signal.name}...")
        self.logger.info("shutting down runner...")
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

        if not self.nc.is_closed:
            self.logger.info("closing nats connection")
            await self.nc.close()

        loop.stop()

    async def run(self) -> None:
        if "default" not in self.response_handlers:
            raise UndefinedDefaultHandlerFunctionError

        if not self.initializer:
            self.initializer = compose_initializer()

        if not self.finalizer:
            self.finalizer = compose_finalizer()

        await self.initializer(self.sdk)

        loop = asyncio.get_event_loop()
        signals = (signal.SIGINT, signal.SIGTERM)
        for s in signals:
            loop.add_signal_handler(
                s,
                lambda s=s: asyncio.create_task(self._shutdown_handler(loop, signal=s)),
            )
        loop.set_exception_handler(self._exception_handler)

        try:
            await self.subscriber.start()
        finally:
            self.logger.info("task runner stopped")

        self.finalizer(self.sdk)
