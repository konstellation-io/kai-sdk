from abc import ABC, abstractmethod
from dataclasses import dataclass
import os

from centralized_configuration.centralized_configuration import CentralizedConfiguration
from google.protobuf.message import Message
from loguru._logger import Logger
from messaging.messaging import Messaging
from metadata.metadata import Metadata
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from object_store.object_store import ObjectStore
from path_utils.path_utils import PathUtils

from kai_nats_msg_pb2 import KaiNatsMessage


class Messaging(ABC):
    @abstractmethod
    def send_output(self, response: Message, *channel_opt: str) -> None:
        pass

    @abstractmethod
    def send_output_with_request_id(self, response: Message, request_id: str, *channel_opt: str) -> None:
        pass

    @abstractmethod
    def send_any(self, response: Message, *channel_opt: str) -> None:
        pass

    @abstractmethod
    def send_early_reply(self, response: Message, *channel_opt: str) -> None:
        pass

    @abstractmethod
    def send_early_exit(self, response: Message, *channel_opt: str) -> None:
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


class Metadata(ABC):
    @abstractmethod
    def get_process(self) -> str:
        pass

    @abstractmethod
    def get_workflow(self) -> str:
        pass

    @abstractmethod
    def get_product(self) -> str:
        pass

    @abstractmethod
    def get_version(self) -> str:
        pass

    @abstractmethod
    def get_object_store_name(self) -> str:
        pass

    @abstractmethod
    def get_key_value_store_product_name(self) -> str:
        pass

    @abstractmethod
    def get_key_value_store_workflow_name(self) -> str:
        pass

    @abstractmethod
    def get_key_value_store_process_name(self) -> str:
        pass


class ObjectStore(ABC):
    @abstractmethod
    def list(self, regexp: str) -> list:
        pass

    @abstractmethod
    def get(self, key: str) -> bytes:
        pass

    @abstractmethod
    def save(self, key: str, value: bytes) -> None:
        pass

    @abstractmethod
    def delete(self, key: str) -> None:
        pass

    @abstractmethod
    def purge(self, regexp: str) -> None:
        pass


class CentralizedConfig(ABC):
    @abstractmethod
    def get_config(self, key: str, scope: str) -> str:
        pass

    @abstractmethod
    def set_config(self, key: str, value: str, scope: str) -> None:
        pass

    @abstractmethod
    def delete_config(self, key: str, scope: str) -> None:
        pass


class PathUtils(ABC):
    @abstractmethod
    def get_base_path(self) -> str:
        pass

    @abstractmethod
    def compose_path(self, *relative_path: str) -> str:
        pass


class Measurements(ABC):
    pass


class Storage(ABC):
    pass


@dataclass
class KaiSDK:
    _nats: NatsClient
    _jetstream: JetStreamContext
    _request_message: KaiNatsMessage

    logger: Logger
    metadata: Metadata
    messaging: Messaging
    object_store: ObjectStore
    centralized_config: CentralizedConfiguration
    path_utils: PathUtils
    measurements: Measurements
    storage: Storage

    def __post_init__(self):
        self.logger = self.logger if self.logger else None
        self.path_utils = self.path_utils if self.path_utils else None
        self.metadata = self.metadata if self.metadata else None
        self.messaging = self.messaging if self.messaging else None
        self.object_store = self.object_store if self.object_store else None
        self.centralized_config = self.centralized_config if self.centralized_config else None
        self.measurements = self.measurements if self.measurements else None
        self.storage = self.storage if self.storage else None

        # try:
        #     from centralized_configuration.centralized_configuration import NewCentralizedConfiguration

        #     self.centralized_config = NewCentralizedConfiguration(self.logger, self.jetstream)
        # except Exception as err:
        #     self.logger.error(f"Error initializing Centralized Configuration: {err}")
        #     os.exit(1)

        # try:
        #     from object_store.object_store import NewObjectStore

        #     self.object_store = NewObjectStore(self.logger, self.jetstream)
        # except Exception as err:
        #     self.logger.error(f"Error initializing Object Store: {err}")
        #     os.exit(1)

    def get_request_id(self):
        return self.request_message.RequestId if self.request_message else None

    # @staticmethod
    # def shallow_copy_with_request(sdk, request_msg):
    #     h_sdk = KaiSDK(sdk.nats, sdk.jetstream, request_msg, sdk.logger)
    #     h_sdk.messaging = sdk.messaging if sdk.messaging else None
    #     # Other assignments can be added for path_utils, metadata, measurements, and storage.

    #     return h_sdk
