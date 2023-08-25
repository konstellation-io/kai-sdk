from __future__ import annotations

from dataclasses import dataclass, field
from queue import Queue
from threading import Event, Thread
from typing import Any, Callable, Optional

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
    response_handler: ResponseHandler = field(init=False)
    response_channels: dict[str, Queue] = field(default_factory=dict)
    initializer: Optional[Initializer] = None
    runner: RunnerFunc = field(init=False)
    subscriber: TriggerSubscriber = field(init=False)
    finalizer: Optional[Finalizer] = None
    runner_thread_shutdown_event: Event = field(default_factory=Event)

    def __post_init__(self) -> None:
        self.sdk = KaiSDK(nc=self.nc, js=self.js, logger=self.logger)
        self.subscriber = TriggerSubscriber(self)

    def with_initializer(self, initializer: Initializer) -> TriggerRunner:
        self.initializer = compose_initializer(initializer)
        return self

    def with_runner(self, runner: RunnerFunc) -> TriggerRunner:
        self.runner = compose_runner(self, runner)
        return self

    def with_finalizer(self, finalizer: Finalizer) -> TriggerRunner:
        self.finalizer = compose_finalizer(finalizer)
        return self

    def get_response_channel(self, request_id: str) -> Any:
        if request_id not in self.response_channels:
            self.response_channels[request_id] = Queue(maxsize=1)
        return self.response_channels[request_id].get()

    async def run(self) -> None:
        if self.runner is None:
            raise UndefinedRunnerFunctionError

        if not self.initializer:
            self.initializer = compose_initializer()

        self.response_handler = get_response_handler(self.response_channels)

        if not self.finalizer:
            self.finalizer = compose_finalizer()

        initializer_func = self.initializer(self.sdk)
        await initializer_func

        runner_thread = Thread(target=self.runner, args=(self, self.sdk))
        runner_thread.start()

        subscriber_thread = Thread(target=self.subscriber.start, args=())
        subscriber_thread.start()

        runner_thread.join()
        subscriber_thread.join()

        self.finalizer(self.sdk)


RunnerFunc = Callable[[TriggerRunner, KaiSDK], None]
