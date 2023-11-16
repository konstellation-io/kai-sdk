from __future__ import annotations

import asyncio
import sys
from abc import ABC
from dataclasses import dataclass, field
from typing import Optional

import loguru
from loguru import logger
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext

from sdk.centralized_config.centralized_config import CentralizedConfig, CentralizedConfigABC
from sdk.ephemeral_storage.ephemeral_storage import EphemeralStorage, EphemeralStorageABC
from sdk.kai_nats_msg_pb2 import KaiNatsMessage
from sdk.messaging.messaging import Messaging, MessagingABC
from sdk.metadata.metadata import Metadata, MetadataABC
from sdk.path_utils.path_utils import PathUtils, PathUtilsABC
from sdk.persistent_storage.persistent_storage import PersistentStorage, PersistentStorageABC

LOGGER_FORMAT = (
    "<green>{time:YYYY-MM-DDTHH:mm:ss.SSS}Z</green> "
    "<cyan>{level}</cyan> {extra[context]} <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> "
    "<level>{message}</level> <level>{extra[request_id]}</level>"
)

@dataclass
class MeasurementsABC(ABC):
    pass


@dataclass
class Storage:
    persistent: PersistentStorageABC
    ephemeral: EphemeralStorageABC = field(default=None)


@dataclass
class KaiSDK:
    nc: NatsClient
    js: JetStreamContext
    logger: Optional[loguru.Logger] = None
    request_msg: KaiNatsMessage = field(init=False, default=None)
    metadata: MetadataABC = field(init=False)
    messaging: MessagingABC = field(init=False)
    centralized_config: CentralizedConfigABC = field(init=False)
    path_utils: PathUtilsABC = field(init=False)
    measurements: MeasurementsABC = field(init=False)
    storage: Storage = field(init=False)

    def __post_init__(self) -> None:
        self.metadata = Metadata()

        if not self.logger:
            self._initialize_logger()
        else:
            self.logger.configure(extra={"context": "[KAI SDK]", "request_id": "{}"})

        self.centralized_config = CentralizedConfig(js=self.js)
        self.messaging = Messaging(nc=self.nc, js=self.js)
        self.path_utils = PathUtils()
        self.measurements = MeasurementsABC()
        self.storage = Storage(PersistentStorage(), EphemeralStorage(js=self.js))

    async def initialize(self) -> None:
        try:
            await self.storage.ephemeral.initialize()
        except Exception as e:
            assert self.logger is not None
            self.logger.error(f"error initializing object store: {e}")
            asyncio.get_event_loop().stop()
            sys.exit(1)

        try:
            await self.centralized_config.initialize()
        except Exception as e:
            assert self.logger is not None
            self.logger.error(f"error initializing centralized configuration: {e}")
            asyncio.get_event_loop().stop()
            sys.exit(1)

    def get_request_id(self) -> str | None:
        request_msg = getattr(self, "request_msg", None)
        return self.request_msg.request_id if request_msg else None

    def _initialize_logger(self) -> None:
        logger.remove()  # Remove the pre-configured handler
        logger.add(
            sys.stdout,
            colorize=True,
            format=LOGGER_FORMAT,
            backtrace=True,
            diagnose=True,
            level="DEBUG",
        )
        logger.configure(extra={"context": "[UNKNOWN]", "request_id": "{}"})

        self.logger = logger.bind(context="[KAI SDK]")
        self.logger.debug("logger initialized")

    def set_request_msg(self, request_msg: KaiNatsMessage) -> None:
        self.request_msg = request_msg
        assert isinstance(self.messaging, Messaging)
        self.messaging.request_msg = request_msg
