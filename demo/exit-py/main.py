import asyncio
import io

from google.protobuf.any_pb2 import Any
from google.protobuf.wrappers_pb2 import StringValue
from runner.runner import Runner
from sdk import kai_sdk as sdk


async def initializer(kai_sdk: sdk.KaiSDK):
    logger = kai_sdk.logger.bind(context="[EXIT INITIALIZER]")
    logger.info("starting example...")
    value, _ = await kai_sdk.centralized_config.get_config("test")
    if value is None:
        logger.info("config value not found")
    else:
        logger.info(f"config value retrieved! {value}")

    value, _ = await kai_sdk.storage.ephemeral.get("test")
    if value is None:
        logger.info("ephemeral storage value not found")
    else:
        logger.info(f"ephemeral storage value retrieved! {value.decode('utf-8')}")

    object_id = "some-object"
    object_value = io.BytesIO(b"some-value")
    version_id = kai_sdk.storage.persistent.save(object_id, object_value)
    logger.info(
        f"persistent storage value set! {object_id} {object_value.getvalue()!s} with version {version_id}"
    )
    obj, _ = kai_sdk.storage.persistent.get("some-object", version_id)
    logger.info(f"persistent storage value retrieved! {object_id} {obj!s}")
    objs = kai_sdk.storage.persistent.list()
    logger.info(f"persistent storage value listed! {objs!s}")
    objs = kai_sdk.storage.persistent.list_versions(object_id)
    logger.info(
        f"persistent storage value listed for objects with prefix {object_id}! {objs!s}"
    )
    kai_sdk.storage.persistent.delete(object_id, version_id)
    logger.info(f"persistent storage value deleted! {object_id}")
    kai_sdk.storage.persistent.get("some-object", version_id)


async def handler(kai_sdk: sdk.KaiSDK, response: Any):
    logger = kai_sdk.logger.bind(context="[EXIT HANDLER]")
    string_value = StringValue()

    response.Unpack(string_value)
    message = string_value.value
    logger.info(f"received message {message}")

    composed_string = f"{message}, processed by the exit process!"
    logger.info(f"sending message {composed_string}")
    await kai_sdk.messaging.send_output(StringValue(value=composed_string))


async def preprocessor(kai_sdk: sdk.KaiSDK, response: Any):
    logger = kai_sdk.logger.bind(context="[EXIT PREPROCESSOR]")
    logger.info("I am an async preprocessor")
    await asyncio.sleep(0.0000001)
    logger.info("I am an async preprocessor, after sleep")


def postprocessor(kai_sdk: sdk.KaiSDK, response: Any):
    logger = kai_sdk.logger.bind(context="[EXIT POSTPROCESSOR]")
    logger.info("I am a sync postprocessor")


async def finalizer(kai_sdk: sdk.KaiSDK):
    logger = kai_sdk.logger.bind(context="[EXIT FINALIZER]")
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
