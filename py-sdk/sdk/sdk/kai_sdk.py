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
from sdk.kai_nats_msg_pb2 import KaiNatsMessage
from sdk.messaging.messaging import Messaging, MessagingABC
from sdk.metadata.metadata import Metadata, MetadataABC
from sdk.object_store.object_store import ObjectStore, ObjectStoreABC
from sdk.path_utils.path_utils import PathUtils, PathUtilsABC


@dataclass
class MeasurementsABC(ABC):
    pass


@dataclass
class StorageABC(ABC):
    pass


@dataclass
class KaiSDK:
    nc: NatsClient
    js: JetStreamContext
    logger: Optional[loguru.Logger] = None
    request_msg: KaiNatsMessage = field(init=False, default=None)
    metadata: MetadataABC = field(init=False)
    messaging: MessagingABC = field(init=False)
    object_store: ObjectStoreABC = field(init=False)
    centralized_config: CentralizedConfigABC = field(init=False)
    path_utils: PathUtilsABC = field(init=False)
    measurements: MeasurementsABC = field(init=False)
    storage: StorageABC = field(init=False)

    def __post_init__(self) -> None:
        self.metadata = Metadata()

        if not self.logger:
            self._initialize_logger()
        else:
            product_id = self.metadata.get_product()
            version_id = self.metadata.get_version()
            workflow_id = self.metadata.get_workflow()
            process_id = self.metadata.get_process()
            metadata_info = f"{product_id=} {version_id=} {workflow_id=} {process_id=}"
            self.logger = self.logger.configure(extra={"context": "[KAI SDK]", "metadata_info": metadata_info})

        self.centralized_config = CentralizedConfig(js=self.js)
        self.messaging = Messaging(nc=self.nc, js=self.js)
        self.object_store = ObjectStore(js=self.js)
        self.path_utils = PathUtils()
        self.measurements = MeasurementsABC()
        self.storage = StorageABC()

    async def initialize(self) -> None:
        try:
            await self.object_store.initialize()
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
            format=(
                "<green>{time:YYYY-MM-DD HH:mm:ss.SSS}</green> | "
                "<cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> | "
                "{extra[context]}: <level>{message}</level> - {extra[metadata_info]}"
            ),
            backtrace=True,
            diagnose=True,
        )
        product_id = self.metadata.get_product()
        version_id = self.metadata.get_version()
        workflow_id = self.metadata.get_workflow()
        process_id = self.metadata.get_process()
        metadata_info = f"{product_id=} {version_id=} {workflow_id=} {process_id=}"
        logger.configure(extra={"context": "[UNKNOW]", "metadata_info": metadata_info})

        self.logger = logger.bind(context="[KAI SDK]")
        self.logger.debug("logger initialized")

    def set_request_msg(self, request_msg: KaiNatsMessage) -> None:
        self.request_msg = request_msg
        assert isinstance(self.messaging, Messaging)
        self.messaging.request_msg = request_msg
