import os
import sys
from abc import ABC, abstractmethod
from dataclasses import dataclass
from typing import Optional

from centralized_config.centralized_config import CentralizedConfig
from google.protobuf.message import Message
from loguru import logger
from loguru._logger import Logger
from messaging.messaging import Messaging
from metadata.metadata import Metadata
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from object_store.object_store import ObjectStore
from path_utils.path_utils import PathUtils

from kai_nats_msg_pb2 import KaiNatsMessage

logger.remove()
logger.add(
    sys.stdout,
    colorize=True,
    format="<green>{time}</green> <level>{extra[context]} {message}</level>",
    backtrace=True,
    diagnose=True,
)


@dataclass
class Messaging(ABC):
    @abstractmethod
    def send_output(self, response: Message, chan: Optional[str]):
        pass

    @abstractmethod
    def send_output_with_request_id(self, response: Message, request_id: str, chan: Optional[str]):
        pass

    @abstractmethod
    def send_any(self, response: Message, chan: Optional[str]):
        pass

    @abstractmethod
    def send_early_reply(self, response: Message, chan: Optional[str]):
        pass

    @abstractmethod
    def send_early_exit(self, response: Message, chan: Optional[str]):
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
class Metadata(ABC):
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
class ObjectStore(ABC):
    @abstractmethod
    def initialize(self):
        pass

    @abstractmethod
    def list(self, regexp: Optional[str]) -> list[str]:
        pass

    @abstractmethod
    def get(self, key: str) -> bytes:
        pass

    @abstractmethod
    def save(self, key: str, value: bytes):
        pass

    @abstractmethod
    def delete(self, key: str):
        pass

    @abstractmethod
    def purge(self, regexp: Optional[str]):
        pass


@dataclass
class CentralizedConfig(ABC):
    @abstractmethod
    def initialize(self):
        pass

    @abstractmethod
    def get_config(self, key: str, scope: str) -> str:
        pass

    @abstractmethod
    def set_config(self, key: str, value: str, scope: str):
        pass

    @abstractmethod
    def delete_config(self, key: str, scope: str):
        pass


@dataclass
class PathUtils(ABC):
    @staticmethod
    @abstractmethod
    def get_base_path(self) -> str:
        pass

    @staticmethod
    @abstractmethod
    def compose_path(self, *relative_path: tuple[str]) -> str:
        pass


@dataclass
class Measurements(ABC):
    pass


@dataclass
class Storage(ABC):
    pass


@dataclass
class KaiSDK:
    nc: NatsClient
    js: JetStreamContext
    req_msg: KaiNatsMessage

    logger: Logger = logger.bind(context="[KAI SDK]")
    metadata: Metadata = Metadata()
    messaging: Messaging = None
    object_store: Optional[ObjectStore] = None
    centralized_config: CentralizedConfig = None
    path_utils: PathUtils = PathUtils()
    measurements: Measurements = None
    storage: Storage = None

    def __post_init__(self):
        self.centralized_config = CentralizedConfig(js=self.js)
        self.object_store = ObjectStore(js=self.js)

    async def initialize(self):
        try:
            await self.object_store.initialize()
        except Exception as e:
            self.logger.error(f"error initializing object store: {e}")
            os.exit(1)

        try:
            await self.centralized_config.initialize()
        except Exception as e:
            self.logger.error(f"error initializing centralized configuration: {e}")
            os.exit(1)

        self.messaging = Messaging(nc=self.nc, js=self.js, req_msg=self.req_msg)

    def get_request_id(self):
        return self.req_msg.RequestId if self.req_msg else None
