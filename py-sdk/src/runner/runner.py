from dataclasses import dataclass
from loguru import logger
from loguru._logger import Logger
import sys
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v


logger.remove()  # Remove the pre-configured handler
logger.add(
    sys.stdout,
    colorize=True,
    format="<green>{time}</green> <level>{extra[context]} {message}</level>",
    backtrace=True,
    diagnose=True,
)


@dataclass
class Runner:
    logger: Logger = logger.bind(context="[KAI RUNNER]")
    nc: NatsClient
    js: JetStreamContext

    def __post_init__(self):
        self.logger.info("Runner initialized")

    async def initialize(self):
        self.logger.info("Runner initializing")
        self.logger.info("Runner initialized")