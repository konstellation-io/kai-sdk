import gzip
from abc import ABC, abstractmethod
from dataclasses import dataclass

from loguru import logger
from loguru._logger import Logger
from messaging.exceptions import FailedGettingMaxMessageSizeError
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v

GZIP_HEADER = b"\x1f\x8b"
GZIP_BEST_COMPRESSION = 9


@dataclass
class MessagingUtils(ABC):
    @abstractmethod
    async def get_max_message_size(self) -> int | str:
        pass


@dataclass
class MessagingUtils:
    js: JetStreamContext
    nc: NatsClient
    logger: Logger = logger.bind(context="[MESSAGING UTILS]")

    async def get_max_message_size(self) -> int | str:
        try:
            stream_info = await self.js.stream_info(v.get("nats.stream"))
        except Exception as e:
            self.logger.warning(f"failed getting stream info: {e}")
            raise FailedGettingMaxMessageSizeError(error=e)

        stream_max_size = int(stream_info.config.max_msg_size)
        server_max_size = self.nc.max_payload()

        if stream_max_size == -1:
            return server_max_size

        if stream_max_size < server_max_size:
            return stream_max_size

        return server_max_size


def size_in_mb(size: int) -> str:
    return f"{size / 1024 / 1024:.2f} MB"


def size_in_kb(size: int) -> str:
    return f"{size / 1024:.2f} KB"


def compress(payload: bytes) -> bytes:
    return gzip.compress(payload, compresslevel=GZIP_BEST_COMPRESSION)


def decompress(payload: bytes) -> bytes:
    return gzip.decompress(payload)


def is_compressed(payload: bytes) -> bool:
    return payload.startswith(GZIP_HEADER)
