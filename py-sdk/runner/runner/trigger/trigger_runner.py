from __future__ import annotations

from dataclasses import dataclass, field
from typing import Callable, Optional

import loguru
from google.protobuf.any_pb2 import Any
from loguru import logger
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext

from runner.common.common import Finalizer, Initializer
from runner.trigger.exceptions import UndefinedRunnerFunctionError
from runner.trigger.subscriber import start_subscriber
from runner.trigger.utils import compose_finalizer, compose_initializer, compose_runner, get_response_handler
from sdk.kai_sdk import KaiSDK

ResponseHandler = Callable[[KaiSDK, Any], Optional[Exception]] # TODO revisit optional exception


@dataclass
class TriggerRunner:
    sdk: KaiSDK = field(init=False)
    nc: NatsClient
    js: JetStreamContext
    logger: loguru.Logger = logger.bind(context="[TRIGGER]")
    response_handler: ResponseHandler = None
    response_channels: dict[str, NatsClient] = field(default_factory=dict)
    initializer: Initializer = None
    runner: RunnerFunc = None
    finalizer: Finalizer = None

    def __post_init__(self):
        self.sdk = KaiSDK(nc=self.nc, js=self.js, logger=self.logger)

    def with_initializer(self, initializer: Initializer):
        self.initializer = compose_initializer(initializer)
        return self

    def with_runner(self, runner: RunnerFunc):
        self.runner = compose_runner(self, runner)
        return self

    def with_finalizer(self, finalizer: Finalizer):
        self.finalizer = compose_finalizer(finalizer)
        return self

    def get_response_channel(self, request_id: str) -> NatsClient:
        if request_id not in self.response_channels:
            self.response_channels[request_id] = self.nc.new_inbox()
        return self.response_channels[request_id]

    def run(self):
        if not self.runner:
            raise UndefinedRunnerFunctionError()

        if not self.initializer:
            self.initializer = compose_initializer(None)

        self.response_handler = get_response_handler(self.response_channels)

        if not self.finalizer:
            self.finalizer = compose_finalizer(None)

        self.initializer(self.sdk)  # await?

        self.runner(self.sdk)  # routine

        start_subscriber()  # routine

        # wait?

        self.finalizer(self.sdk)


RunnerFunc = Callable[[TriggerRunner, KaiSDK], None]
