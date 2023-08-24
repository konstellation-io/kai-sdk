from __future__ import annotations

import asyncio
import sys
from dataclasses import dataclass, field
from datetime import timedelta
from signal import SIGINT, SIGTERM
from threading import Event
from typing import TYPE_CHECKING

import loguru
from nats.aio.subscription import Msg
from nats.js.api import ConsumerConfig, DeliverPolicy
from nats.js.client import JetStreamContext
from vyper import v

from sdk.messaging.messaging_utils import is_compressed, uncompress

if TYPE_CHECKING:
    from runner.trigger.trigger_runner import TriggerRunner

from asyncio import AbstractEventLoop

from runner.trigger.exceptions import HandlerError, NewRequestMsgError, NotValidProtobuf, UndefinedResponseHandlerError
from sdk.kai_nats_msg_pb2 import KaiNatsMessage

ACK_TIME = 22  # hours


@dataclass
class TriggerSubscriber:
    trigger_runner: "TriggerRunner"
    logger: loguru.Logger = field(init=False)

    def __post_init__(self) -> None:
        self.logger = self.trigger_runner.logger.bind(context="[SUBSCRIBER]")

    async def start(self) -> None:
        input_subjects = v.get("nats.inputs")
        subscriptions: list[JetStreamContext.PushSubscription] = []
        process = self.trigger_runner.sdk.metadata.get_process().replace(".", "-").replace(" ", "-")

        ack_wait_time = timedelta(hours=ACK_TIME)
        subscriber_thread_shutdown_event = Event()
        for _, subject in input_subjects:
            subject_ = subject.replace(".", "-")
            consumer_name = f"{subject_}_{process}"

            self.logger.info(f"subscribing to {subject} from queue group {consumer_name}")
            try:
                sub = await self.trigger_runner.js.subscribe(
                    subject=subject,
                    queue=consumer_name,
                    cb=self._process_message,
                    deliver_policy=DeliverPolicy.NEW,
                    durable=consumer_name,
                    manual_ack=True,
                    config=ConsumerConfig(ack_wait=ack_wait_time.total_seconds()),
                )
            except Exception as e:
                self.logger.error(f"error subscribing to the NATS subject {subject}: {e}")
                subscriber_thread_shutdown_event.set()
                asyncio.get_event_loop().stop()
                sys.exit(1)

            subscriptions.append(sub)
            self.logger.info(f"listening to {subject} from queue group {consumer_name}")

        def _shutdown_handler(loop: AbstractEventLoop):
            loop.create_task(_shutdown_handler_coro())

        async def _shutdown_handler_coro() -> None:
            self.logger.info("shutting signal received")

            for sub in subscriptions:
                self.logger.info(f"unsubscribing from subject {sub.subject}")

                try:
                    await sub.unsubscribe()
                except Exception as e:
                    self.logger.error(f"error unsubscribing from the NATS subject {sub.subject}: {e}")
                    subscriber_thread_shutdown_event.set()
                    asyncio.get_event_loop().stop()
                    sys.exit(1)

            subscriber_thread_shutdown_event.set()

        loop = asyncio.get_event_loop()
        loop.add_signal_handler(SIGINT, _shutdown_handler(loop))
        loop.add_signal_handler(SIGTERM, _shutdown_handler(loop))

        subscriber_thread_shutdown_event.wait()
        self.logger.info("subscriber shutdown")

    async def _process_message(self, msg: Msg) -> None:
        self.logger.info("new message received")

        try:
            request_msg = self._new_request_msg(msg.data)
            self.trigger_runner.sdk.set_request_message(request_msg)
        except Exception as e:
            await self._process_runner_error(msg, NotValidProtobuf(msg.subject, error=e), request_msg.request_id)
            return

        self.logger.info(f"processing message with request_id {request_msg.request_id} and subject {msg.subject}")

        if self.trigger_runner.response_handler is None:
            await self._process_runner_error(msg, UndefinedResponseHandlerError(), request_msg.request_id)
            return

        try:
            await self.trigger_runner.response_handler(self.trigger_runner.sdk, request_msg.payload)
        except Exception as e:
            from_node = request_msg.from_node
            to_node = self.trigger_runner.sdk.metadata.get_process()
            await self._process_runner_error(msg, HandlerError(from_node, to_node, error=e), request_msg.request_id)
            return

        try:
            await msg.ack()
        except Exception as e:
            self.logger.error(f"error acknowledging message: {e}")

    async def _process_runner_error(self, msg: Msg, error: Exception, request_id: str) -> None:
        error_msg = str(error)
        self.logger.info(f"publishing error message {error_msg}")

        try:
            await msg.ack()
        except Exception as e:
            self.logger.error(f"error acknowledging message: {e}")

        await self.trigger_runner.sdk.messaging.send_error(error_msg, request_id)

    def _new_request_msg(self, data: bytes) -> KaiNatsMessage:
        request_msg = KaiNatsMessage()

        if is_compressed(data):
            try:
                data = uncompress(data)
            except Exception as e:
                error = NewRequestMsgError(error=e)
                self.logger.error(f"{error}")
                raise error

        try:
            request_msg.ParseFromString(data)  # deserialize from bytes
        except Exception as e:
            error = NewRequestMsgError(error=e)
            self.logger.error(f"{error}")
            raise error

        return request_msg
