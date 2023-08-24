from __future__ import annotations

from dataclasses import dataclass, field
from typing import Optional

import loguru
from loguru import logger
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v

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

Preprocessor = Handler
Postprocessor = Handler


@dataclass
class TaskRunner:
    sdk: KaiSDK = field(init=False)
    nc: NatsClient
    js: JetStreamContext
    logger: loguru.Logger = logger.bind(context="[TASK]")
    response_handlers: dict[str, Handler] = field(default_factory=dict)
    initializer: Optional[Initializer] = None
    preprocessor: Preprocessor = field(init=False)
    postprocessor: Postprocessor = field(init=False)
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

    async def run(self) -> None:
        if "default" not in self.response_handlers:
            raise UndefinedDefaultHandlerFunctionError()

        if not self.initializer:
            self.initializer = compose_initializer()

        if not self.finalizer:
            self.finalizer = compose_finalizer()

        initializer_func = self.initializer(self.sdk)
        await initializer_func

        await self.subscriber.start()

        self.finalizer(self.sdk)
