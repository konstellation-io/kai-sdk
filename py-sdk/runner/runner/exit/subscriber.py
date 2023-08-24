from __future__ import annotations

import asyncio
import sys
from dataclasses import dataclass, field
from datetime import timedelta
from signal import SIGINT, SIGTERM, signal
from threading import Event
from typing import TYPE_CHECKING

import loguru
from nats.aio.subscription import Msg
from nats.js.api import ConsumerConfig, DeliverPolicy
from nats.js.client import JetStreamContext
from vyper import v

from sdk.messaging.messaging_utils import is_compressed, uncompress

if TYPE_CHECKING:
    from runner.exit.exit_runner import ExitRunner, Handler

from runner.exit.exceptions import HandlerError, NewRequestMsgError, NotValidProtobuf, UndefinedHandlerFunctionError
from sdk.kai_nats_msg_pb2 import KaiNatsMessage
from sdk.metadata.metadata import Metadata

ACK_TIME = 22  # hours


@dataclass
class ExitSubscriber:
    exit_runner: "ExitRunner"
    logger: loguru.Logger = field(init=False)

    def __post_init__(self) -> None:
        self.logger = self.exit_runner.logger.bind(context="[SUBSCRIBER]")

    async def start(self) -> None:
        input_subjects = v.get("nats.inputs")
        subscriptions: list[JetStreamContext.PushSubscription] = []
        process = self.exit_runner.sdk.metadata.get_process().replace(".", "-").replace(" ", "-")

        for _, subject in input_subjects:
            subject_ = subject.replace(".", "-")
            consumer_name = f"{subject_}_{process}"

            self.logger.info(f"subscribing to {subject} from queue group {consumer_name}")

            ack_wait_time = timedelta(hours=ACK_TIME)

            subscriber_thread_shutdown_event = Event()
            try:
                sub = await self.exit_runner.js.subscribe(
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

        async def _shutdown_handler(sig, frame) -> None:
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

        signal(SIGINT, asyncio.create_task(_shutdown_handler))
        signal(SIGTERM, asyncio.create_task(_shutdown_handler))

        subscriber_thread_shutdown_event.wait()
        self.logger.info("subscriber shutdown")

    async def _process_message(self, msg: Msg) -> None:
        self.logger.info("new message received")

        try:
            request_msg = self._new_request_msg(msg.data)
            self.exit_runner.sdk.set_request_message(request_msg)
        except Exception as e:
            await self._process_runner_error(msg, NotValidProtobuf(msg.subject, error=e), request_msg.request_id)
            return

        self.logger.info(f"processing message with request_id {request_msg.request_id} and subject {msg.subject}")

        from_node = request_msg.from_node
        handler = self._get_response_handler(from_node.lower())
        to_node = self.exit_runner.sdk.metadata.get_process()

        if handler is None:
            await self._process_runner_error(msg, UndefinedHandlerFunctionError(from_node), request_msg.request_id)
            return

        try:
            if self.exit_runner.preprocessor is not None:
                self.exit_runner.preprocessor(self.exit_runner.sdk, request_msg.payload)
        except Exception as e:
            await self._process_runner_error(
                msg, HandlerError(from_node, to_node, error=e, type="handler preprocessor"), request_msg.request_id
            )
            return

        try:
            handler(self.exit_runner.sdk, request_msg.payload)
        except Exception as e:
            await self._process_runner_error(msg, HandlerError(from_node, to_node, error=e), request_msg.request_id)
            return

        try:
            if self.exit_runner.postprocessor is not None:
                self.exit_runner.postprocessor(self.exit_runner.sdk, request_msg.payload)
        except Exception as e:
            await self._process_runner_error(
                msg, HandlerError(from_node, to_node, error=e, type="handler postprocessor"), request_msg.request_id
            )
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

        await self.exit_runner.sdk.messaging.send_error(error_msg, request_id)

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

    def _get_response_handler(self, subject: str) -> Handler:
        if subject in self.exit_runner.response_handlers:
            return self.exit_runner.response_handlers[subject]

        return self.exit_runner.response_handlers["default"]
