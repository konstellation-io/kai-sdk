from abc import ABC, abstractmethod
from dataclasses import dataclass

from exceptions import FailedGettingMaxMessageSizeError
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v


@dataclass
class MessagingUtils(ABC):
    @abstractmethod
    async def get_max_message_size(self) -> int | str:
        pass


@dataclass
class MessagingUtils:
    js: JetStreamContext
    nc: NatsClient

    async def get_max_message_size(self) -> int | str:
        try:
            stream_info = await self.js.stream_info(v.get("nats.stream"))
        except Exception as e:
            raise FailedGettingMaxMessageSizeError(error=e)

        stream_max_size = int(stream_info.config.max_msg_size)
        server_max_size = self.nc.max_payload()

        if stream_max_size == -1:
            return server_max_size

        if stream_max_size < server_max_size:
            return stream_max_size

        return server_max_size


def size_in_mb(size: int) -> str:
    return f"{size / 1024 / 1024:.1f} MB"
