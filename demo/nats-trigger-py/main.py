import asyncio
import uuid
from datetime import timedelta

from google.protobuf.wrappers_pb2 import StringValue
from nats.aio.client import Client as NatsClient
from nats.aio.msg import Msg
from nats.js.api import ConsumerConfig, DeliverPolicy
from nats.js.client import JetStreamContext
from runner.runner import Runner
from runner.trigger import trigger_runner
from sdk import kai_sdk as sdk


async def initializer(sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[NATS SUBSCRIBER INITIALIZER]")
    logger.info("starting example...")


async def nats_subscriber_runner(
    trigger_runner: trigger_runner.TriggerRunner, sdk: sdk.KaiSDK
):
    logger = sdk.logger.bind(context="[NATS SUBSCRIBER RUNNER]")
    logger.info("executing example...")

    nc: NatsClient = NatsClient()
    js: JetStreamContext = nc.jetstream()
    await nc.connect("nats://localhost:4222")
    ack_time = timedelta(hours=22).total_seconds()

    async def process_message(msg: Msg):
        logger = sdk.logger.bind(context="[NATS SUBSCRIBER CALLBACK]")
        logger.info("processing message...")

        request_id = str(uuid.uuid4())
        await sdk.messaging.send_output_with_request_id(
            StringValue(value="Hi there, I'm a NATS subscriber!"), request_id
        )

        response_channel = await trigger_runner.get_response_channel(request_id)

        logger.info(f"message received: {response_channel}")

        try:
            await msg.ack()
        except Exception as e:
            logger.error(f"error acknowledging message: {e}")

    logger.info("subscribing to nats-trigger...")
    try:
        await js.subscribe(
            subject="demo-trigger",
            queue="demo-trigger-queue",
            durable="demo-trigger",
            cb=process_message,
            config=ConsumerConfig(deliver_policy=DeliverPolicy.NEW, ack_wait=ack_time),
            manual_ack=True,
        )
    except Exception as e:
        logger.error(f"error subscribing to the NATS subject demo-trigger: {e}")
        return


def finalizer(sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[NATS SUBSCRIBER FINALIZER]")
    logger.info("finalizing example...")


async def init():
    runner = await Runner().initialize()
    await runner.trigger_runner().with_initializer(initializer).with_runner(
        nats_subscriber_runner
    ).with_finalizer(finalizer).run()


if __name__ == "__main__":
    loop = asyncio.new_event_loop()
    loop.run_until_complete(init())
    loop.run_forever()
    loop.close()
