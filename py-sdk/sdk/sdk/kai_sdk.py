from __future__ import annotations

import sys
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import Optional

import loguru
from google.protobuf.any_pb2 import Any
from google.protobuf.message import Message
from loguru import logger
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext

from sdk.centralized_config.centralized_config import CentralizedConfig
from sdk.centralized_config.constants import Scope
from sdk.kai_nats_msg_pb2 import KaiNatsMessage
from sdk.messaging.messaging import Messaging
from sdk.metadata.metadata import Metadata
from sdk.object_store.object_store import ObjectStore
from sdk.path_utils.path_utils import PathUtils
import asyncio

@dataclass
class MessagingABC(ABC):
    @abstractmethod
    async def send_output(self, response: Message, chan: Optional[str]):
        pass

    @abstractmethod
    async def send_output_with_request_id(self, response: Message, request_id: str, chan: Optional[str]):
        pass

    @abstractmethod
    async def send_any(self, response: Any, chan: Optional[str]):
        pass

    @abstractmethod
    async def send_any_with_request_id(self, response: Any, request_id: str, chan: Optional[str]):
        pass

    @abstractmethod
    async def send_error(self, error: str, request_id: str):
        pass

    @abstractmethod
    async def send_early_reply(self, response: Message, chan: Optional[str]):
        pass

    @abstractmethod
    async def send_early_exit(self, response: Message, chan: Optional[str]):
        pass

    @abstractmethod
    def is_message_ok(self) -> bool:
        pass

    @abstractmethod
    def is_message_error(self) -> bool:
        pass

    @abstractmethod
    def is_message_early_reply(self) -> bool:
        pass

    @abstractmethod
    def is_message_early_exit(self) -> bool:
        pass


@dataclass
class MetadataABC(ABC):
    @staticmethod
    @abstractmethod
    def get_product(self) -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_workflow(self) -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_process(self) -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_version(self) -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_object_store_name(self) -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_key_value_store_product_name(self) -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_key_value_store_workflow_name(self) -> str:
        pass

    @staticmethod
    @abstractmethod
    def get_key_value_store_process_name(self) -> str:
        pass


@dataclass
class ObjectStoreABC(ABC):
    @abstractmethod
    async def initialize(self) -> Optional[Exception]:
        pass

    @abstractmethod
    async def list(self, regexp: Optional[str]) -> list[str] | Exception:
        pass

    @abstractmethod
    async def get(self, key: str) -> tuple[bytes, bool] | Exception:
        pass

    @abstractmethod
    async def save(self, key: str, payload: bytes) -> Optional[Exception]:
        pass

    @abstractmethod
    async def delete(self, key: str) -> bool | Exception:
        pass

    @abstractmethod
    async def purge(self, regexp: Optional[str]) -> Optional[Exception]:
        pass


@dataclass
class CentralizedConfigABC(ABC):
    @abstractmethod
    async def initialize(self) -> Optional[Exception]:
        pass

    @abstractmethod
    async def get_config(self, key: str, scope: Optional[Scope]) -> tuple[str, bool] | Exception:
        pass

    @abstractmethod
    async def set_config(self, key: str, value: bytes, scope: Optional[Scope]) -> Optional[Exception]:
        pass

    @abstractmethod
    async def delete_config(self, key: str, scope: Optional[Scope]) -> bool | Exception:
        pass


@dataclass
class PathUtilsABC(ABC):
    @staticmethod
    @abstractmethod
    def get_base_path(self) -> str:
        pass

    @staticmethod
    @abstractmethod
    def compose_path(self, *relative_path: str) -> str:
        pass


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
    req_msg: KaiNatsMessage = field(init=False)
    metadata: MetadataABC = field(init=False)
    messaging: MessagingABC = field(init=False)
    object_store: Optional[ObjectStoreABC] = field(init=False)
    centralized_config: CentralizedConfigABC = field(init=False)
    path_utils: PathUtilsABC = field(init=False)
    measurements: MeasurementsABC = field(init=False)
    storage: StorageABC = field(init=False)

    def __post_init__(self):
        if not self.logger:
            self._initialize_logger()
        else:
            self.logger = self.logger.bind(context="[KAI SDK]")

        self.req_msg = None
        self.centralized_config = CentralizedConfig(js=self.js)
        self.metadata = Metadata()
        self.messaging = Messaging(nc=self.nc, js=self.js)
        self.object_store = ObjectStore(js=self.js)
        self.path_utils = PathUtils()
        self.measurements = None
        self.storage = None

    async def initialize(self):
        try:
            await self.object_store.initialize()
        except Exception as e:
            self.logger.error(f"error initializing object store: {e}")
            asyncio.get_event_loop().stop()
            sys.exit(1)

        try:
            await self.centralized_config.initialize()
        except Exception as e:
            self.logger.error(f"error initializing centralized configuration: {e}")
            asyncio.get_event_loop().stop()
            sys.exit(1)

    def get_request_id(self):
        return self.req_msg.request_id if self.req_msg else None

    def _initialize_logger(self):
        logger.remove()  # Remove the pre-configured handler
        logger.add(
            sys.stdout,
            colorize=True,
            format="<green>{time}</green> <level>{extra[context]} {message}</level>",
            backtrace=True,
            diagnose=True,
        )

        self.logger = logger.bind(context="[KAI SDK]")
        self.logger.debug("logger initialized")

    def set_request_message(self, req_msg: KaiNatsMessage):
        self.req_msg = req_msg
        assert isinstance(self.messaging, Messaging)
        self.messaging.req_msg = req_msg
