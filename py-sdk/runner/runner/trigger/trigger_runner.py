from dataclasses import dataclass

from loguru._logger import Logger
from loguru import logger
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from sdk.kai_sdk import KaiSDK
from runner.trigger.exceptions import FailedInitializingConfigError

@dataclass
class TriggerRunner:
    sdk: KaiSDK = None
    nc: NatsClient
    js: JetStreamContext
    logger: Logger = logger.bind(runner="[TRIGGER]")
    initializer: common.Initializer = None
    #runner: Runner = None
    finalizer: common.Finalizer = None

    def __post_init__(self):
        self.sdk = KaiSDK(nc=self.nc, js=self.js, logger=self.logger)

    def with_initializer(self, initializer):
        self.initializer = composeInitializer(initializer)
        return self
    
    def with_runner(self, runner):
        self.runner = composeRunner(runner)
        return self
    
    def with_finalizer(self, finalizer):
        self.finalizer = composeFinalizer(finalizer)
        return self
    

    def run(self):
        if not self.runner:
            raise FailedInitializingConfigError
        
        if not self.initializer:
            self.initializer = composeInitializer(None)

        if not self.finalizer:
            self.finalizer = composeFinalizer(None)

        self.initializer(self.sdk)

        # TODO

        self.finalizer(self.sdk)