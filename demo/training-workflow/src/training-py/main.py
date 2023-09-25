import asyncio

from google.protobuf.any_pb2 import Any
from google.protobuf.wrappers_pb2 import StringValue
from runner.runner import Runner
from sdk import kai_sdk as sdk
from proto.training_pb2 import Splitter, Training


import random
import string

def get_random_string(length: int) -> str:
    letters = string.ascii_lowercase
    return (''.join(random.choice(letters) for _ in range(length)))


async def initializer(sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[TRAINING INITIALIZER]")
    logger.info("starting process...")


async def handler(sdk: sdk.KaiSDK, response: Any):
    logger = sdk.logger.bind(context="[TRAINING HANDLER]")
    logger.info("training task received")
    input_proto = Splitter()

    response.Unpack(input_proto)
    logger.info(f"received message {input_proto}")

    output = Training(
        training_id=input_proto.training_id,
        model_id=get_random_string(random.randint(10, 20))
    )
    logger.info(f"sending message {output}")

    await sdk.messaging.send_output(StringValue(value=output))
    logger.info("training task done")


def finalizer(sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[TRAINING FINALIZER]")
    logger.info("finalizing process...")


async def init():
    runner = await Runner().initialize()
    await runner.task_runner().with_initializer(initializer).with_handler(handler).with_finalizer(finalizer).run()


if __name__ == "__main__":
    loop = asyncio.new_event_loop()
    loop.run_until_complete(init())
    loop.run_forever()
    loop.close()
