from __future__ import annotations

from dataclasses import dataclass, field
from datetime import timedelta
from signal import SIGINT, SIGTERM, signal
import sys

import loguru
from nats.js.api import DeliverPolicy, ConsumerConfig
from nats.js.client import JetStreamContext
from nats.aio.subscription import Msg
import asyncio
from vyper import v
from sdk.messaging.messaging_utils import uncompress, is_compressed
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from runner.trigger.trigger_runner import TriggerRunner
from runner.trigger.exceptions import NewRequestMsgError, UndefinedResponseHandlerError, HandlerError
from runner.kai_nats_msg_pb2 import KaiNatsMessage
from sdk.metadata.metadata import Metadata
ACK_TIME = 22  # hours

@dataclass
class TriggerSubscriber:
    trigger_runner: 'TriggerRunner'
    logger: loguru.Logger = field(init=False)

    def __post_init__(self):
        self.logger = self.trigger_runner.logger.bind(context="[SUBSCRIBER]")

    async def start_subscriber(self):
        input_subjects = v.get("nats.inputs")
        subscriptions: list[JetStreamContext.PushSubscription] = []

        for _, subject in input_subjects:
            subject_ = subject.replace('.', '-')
            process_ = self.trigger_runner.sdk.metadata.get_process().replace('.', '-').replace(' ', '-')
            consumer_name = f"{subject_}_{process_}"

            self.logger.info(f"subscribing to {subject} from queue group {consumer_name}")

            ack_wait_time = timedelta(hours=ACK_TIME)

            try:
                sub = await self.trigger_runner.js.subscribe(
                    subject=subject,
                    queue=consumer_name,
                    cb=self.process_message,
                    deliver_policy=DeliverPolicy.NEW,
                    durable=consumer_name,
                    manual_ack=True,
                    config=ConsumerConfig(ack_wait=ack_wait_time.total_seconds()),
                )
            except Exception as e:
                self.logger.error(f"error subscribing to the NATS subject {subject}: {e}")
                self.trigger_runner.subscriber_thread_shutdown_event.set()
                asyncio.get_event_loop().stop()
                sys.exit(1)

            subscriptions.append(sub)
            self.logger.info(f"listening to {subject} from queue group {consumer_name}")

        async def shutdown_handler(sig, frame):
            self.logger.info("shutting signal received")

            for sub in subscriptions:
                self.logger.info(f"unsubscribing from subject {sub.subject}")

                try:
                    await sub.unsubscribe()
                except Exception as e:
                    self.logger.error(f"error unsubscribing from the NATS subject {sub.subject}: {e}")
                    self.trigger_runner.subscriber_thread_shutdown_event.set()
                    asyncio.get_event_loop().stop()
                    sys.exit(1)

            self.trigger_runner.subscriber_thread_shutdown_event.set()

        signal(SIGINT, asyncio.create_task(shutdown_handler))
        signal(SIGTERM, asyncio.create_task(shutdown_handler))

        self.trigger_runner.subscriber_thread_shutdown_event.wait()
        self.logger.info("subscriber shutdown")

    async def process_message(self, msg: Msg):
        self.logger.info("new message received")

        try:
            request_msg = self.new_request_msg(msg.data)
        except Exception as e:
            self.logger.error(f"error parsing message: {e}")
            await self.process_runner_error(msg, e, request_msg.request_id)
            return
        
        self.logger.info(f"processing message with request_id {request_msg.request_id} and subject {msg.subject}")
        
        if self.trigger_runner.response_handler is None:
            self.logger.error("no response handler defined")
            await self.process_runner_error(msg, UndefinedResponseHandlerError, request_msg.request_id)
            return
        
        self.trigger_runner.sdk.set_request_message(request_msg)

        try:
            await self.trigger_runner.response_handler(self.trigger_runner.sdk, request_msg.payload)
        except Exception as e:
            self.logger.error(f"error executing response handler: {e}")
            assert isinstance(self.trigger_runner.sdk.metadata, Metadata)
            error = HandlerError(request_msg.from_node, self.trigger_runner.sdk.metadata.get_process(), error=e)
            await self.process_runner_error(msg, error, request_msg.request_id)
            return

        try:
            await msg.ack()
        except Exception as e:
            self.logger.error(f"error acknowledging message: {e}")

    async def process_runner_error(self, msg: Msg, error: Exception, request_id: str):
        try:
            await msg.ack()
        except Exception as e:
            self.logger.error(f"error acknowledging message: {e}")

        self.logger.info(f"publishing error message {error}")
        await self.trigger_runner.sdk.messaging.send_error(str(error), request_id)
    
    def new_request_msg(self, data: bytes) -> KaiNatsMessage:
        request_msg = KaiNatsMessage()

        if is_compressed(data):
            try:
                data = uncompress(data)
            except Exception as e:
                self.logger.error(f"error decompressing message: {e}")
                raise NewRequestMsgError(error=e)
            
        try:
            request_msg.ParseFromString(data) # deserialize from bytes
        except Exception as e:
            self.logger.error(f"error parsing message: {e}")
            raise NewRequestMsgError(error=e)
        
        return request_msg
