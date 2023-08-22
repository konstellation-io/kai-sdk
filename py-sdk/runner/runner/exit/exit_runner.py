from __future__ import annotations


from dataclasses import dataclass
from runner.common.common import Handler
from runner.exit.helpers import compose_finalizer, compose_initializer, compose_handler, compose_postprocessor, compose_preprocessor
from runner.exit.subscriber import ExitSubscriber
from sdk.kai_sdk import KaiSDK
from dataclasses import field
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v
import loguru
from typing import Optional
from loguru import logger
from runner.common.common import Finalizer, Initializer
from runner.exit.exceptions import UndefinedDefaultHandlerFunctionError

Preprocessor = Postprocessor = Handler


@dataclass
class ExitRunner:
    sdk: KaiSDK = field(init=False)
    nc: NatsClient
    js: JetStreamContext
    logger: loguru.Logger = logger.bind(context="[EXIT]")
    response_handlers: dict[str, Handler] = field(default_factory=dict)
    initializer: Optional[Initializer] = None
    preprocessor: Preprocessor = field(init=False)
    postprocessor: Postprocessor = field(init=False)
    finalizer: Optional[Finalizer] = None

    def __post_init__(self):
        self.sdk = KaiSDK(nc=self.nc, js=self.js, logger=self.logger)
        self.subscriber = ExitSubscriber(self)

    def with_initializer(self, initializer: Initializer):
        self.initializer = compose_initializer(initializer)
        return self
    
    def with_preprocessor(self, preprocessor: Preprocessor):
        self.preprocessor = compose_preprocessor(preprocessor)
        return self
    
    def with_handler(self, handler: Handler):
        self.response_handlers["default"] = compose_handler(handler)
        return self
    
    def with_custom_handler(self, subject: str, handler: Handler):
        self.response_handlers[subject] = compose_handler(handler)
        return self
    
    def with_postprocessor(self, postprocessor: Postprocessor):
        self.postprocessor = compose_postprocessor(postprocessor)
        return self
    
    def with_finalizer(self, finalizer: Finalizer):
        self.finalizer = compose_finalizer(finalizer)
        return self
    
    async def run(self):
        if "default" not in self.response_handlers:
            raise UndefinedDefaultHandlerFunctionError()
        
        if not self.initializer:
            self.initializer = compose_initializer(None)

        if not self.finalizer:
            self.finalizer = compose_finalizer(None)

        initializer_func = self.initializer(self.sdk)
        await initializer_func

        self.subscriber.start()
        
        self.finalizer(self.sdk)