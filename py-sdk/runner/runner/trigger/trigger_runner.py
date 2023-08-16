from dataclasses import dataclass

from loguru._logger import Logger
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v


@dataclass
class TriggerRunner:
    sdk = None
    nc: NatsClient = NatsClient()
    js: JetStreamContext = None
    logger: Logger = None
