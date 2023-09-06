import asyncio

from google.protobuf.wrappers_pb2 import StringValue
from google.protobuf.any_pb2 import Any

from runner.runner import Runner
from sdk import kai_sdk as sdk
from sdk.centralized_config.centralized_config import Scope


async def initializer(sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[EXIT INITIALIZER]")
    logger.info("starting example...")
    value, _ = await sdk.centralized_config.get_config("test")
    if value is None:
        logger.info("config value not found")
    else:
        logger.info(f"config value retrieved! {value}")

    value, _ = await sdk.object_store.get("test")
    if value is None:
        logger.info("object store value not found")
    else:
        logger.info(f"object store value retrieved! {value.decode('utf-8')}")


async def handler(sdk: sdk.KaiSDK, response: Any):
    logger = sdk.logger.bind(context="[EXIT HANDLER]")
    string_value = StringValue()

    response.Unpack(string_value)
    message = string_value.value
    logger.info(f"received message {message}")

    composed_string = f"{message}, processed by the exit process!"
    logger.info(f"sending message {composed_string}")
    await sdk.messaging.send_output(StringValue(value=composed_string))


async def preprocesor(sdk: sdk.KaiSDK, response: Any):
    logger = sdk.logger.bind(context="[EXIT PREPROCESSOR]")
    logger.info("I am an async preprocessor")
    await asyncio.sleep(0.0000001)
    logger.info("I am an async preprocessor, after sleep")

def postprocesor(sdk: sdk.KaiSDK, response: Any):
    logger = sdk.logger.bind(context="[EXIT POSTPROCESSOR]")
    logger.info("I am a sync postprocessor")

async def finalizer(sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[EXIT FINALIZER]")
    logger.info("finalizing example asyncronously...")
    await asyncio.sleep(0.0000001)
    logger.info("finalizing example asyncronously, after sleep...")

async def init():
    runner = await Runner().initialize()
    await runner.exit_runner().with_initializer(initializer).with_handler(handler).with_finalizer(finalizer).with_preprocessor(preprocesor).with_postprocessor(postprocesor).run()

if __name__ == "__main__":
    loop = asyncio.new_event_loop()
    loop.run_until_complete(init())
    loop.run_forever()
    loop.close()
