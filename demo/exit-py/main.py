import asyncio

from google.protobuf.any_pb2 import Any
from google.protobuf.wrappers_pb2 import StringValue
from runner.runner import Runner
from sdk import kai_sdk as sdk


async def initializer(sdk_: sdk.KaiSDK):
    logger = sdk_.logger.bind(context="[EXIT INITIALIZER]")
    logger.info("starting example...")
    value, _ = await sdk_.centralized_config.get_config("test")
    if value is None:
        logger.info("config value not found")
    else:
        logger.info(f"config value retrieved! {value}")

    value, _ = await sdk_.storage.ephemeral.get("test")
    if value is None:
        logger.info("ephemeral storage value not found")
    else:
        logger.info(f"ephemeral storage value retrieved! {value.decode('utf-8')}")


async def handler(sdk_: sdk.KaiSDK, response: Any):
    logger = sdk_.logger.bind(context="[EXIT HANDLER]")
    string_value = StringValue()

    response.Unpack(string_value)
    message = string_value.value
    logger.info(f"received message {message}")

    composed_string = f"{message}, processed by the exit process!"
    logger.info(f"sending message {composed_string}")
    await sdk_.messaging.send_output(StringValue(value=composed_string))


async def preprocessor(sdk_: sdk.KaiSDK, response: Any):
    logger = sdk_.logger.bind(context="[EXIT PREPROCESSOR]")
    logger.info("I am an async preprocessor")
    await asyncio.sleep(0.0000001)
    logger.info("I am an async preprocessor, after sleep")


def postprocessor(sdk_: sdk.KaiSDK, response: Any):
    logger = sdk_.logger.bind(context="[EXIT POSTPROCESSOR]")
    logger.info("I am a sync postprocessor")


async def finalizer(sdk_: sdk.KaiSDK):
    logger = sdk_.logger.bind(context="[EXIT FINALIZER]")
    logger.info("finalizing example asynchronously...")
    await asyncio.sleep(0.0000001)
    logger.info("finalizing example asynchronously, after sleep...")


async def init():
    runner = await Runner().initialize()
    await runner.exit_runner().with_initializer(initializer).with_handler(
        handler
    ).with_finalizer(finalizer).with_preprocessor(preprocessor).with_postprocessor(
        postprocessor
    ).run()


if __name__ == "__main__":
    loop = asyncio.new_event_loop()
    loop.run_until_complete(init())
    loop.run_forever()
    loop.close()
