import asyncio
import random
import string

from google.protobuf.any_pb2 import Any
from google.protobuf.wrappers_pb2 import StringValue
from proto.training_pb2 import Splitter
from runner.runner import Runner
from sdk import kai_sdk as sdk


def get_random_string(length: int) -> str:
    letters = string.ascii_lowercase
    return "".join(random.choice(letters) for _ in range(length))


async def initializer(sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[SPLITTER INITIALIZER]")
    logger.info("starting process...")


async def handler(sdk: sdk.KaiSDK, response: Any):
    logger = sdk.logger.bind(context="[SPLITTER HANDLER]")
    logger.info("splitting task received")
    input_proto = Splitter()  # TODO

    response.Unpack(input_proto)
    logger.info(f"received repo url {input_proto.repo_url}")

    output = Splitter(
        training_id=get_random_string(random.randint(10, 20)),
        repo_url=input_proto.repo_url,
    )
    logger.info(f"sending message {output}")

    await sdk.messaging.send_output(StringValue(value=output))
    logger.info("splitting task done")


def finalizer(sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[SPLITTER FINALIZER]")
    logger.info("finalizing process...")


async def init():
    runner = await Runner().initialize()
    await runner.task_runner().with_initializer(initializer).with_handler(
        handler
    ).with_finalizer(finalizer).run()


if __name__ == "__main__":
    loop = asyncio.new_event_loop()
    loop.run_until_complete(init())
    loop.run_forever()
    loop.close()
