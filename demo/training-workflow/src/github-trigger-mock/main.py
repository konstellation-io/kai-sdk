import asyncio

from google.protobuf.struct_pb2 import Struct
from runner.runner import Runner
from runner.trigger import trigger_runner
from sdk import kai_sdk as sdk


async def initializer(sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[GITHUB-TRIGGER INITIALIZER]")
    logger.info("starting process...")


async def basic_runner(trigger_runner: trigger_runner, sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[GITHUB-TRIGGER HANDLER]")
    logger.info("github repo event received")

    output = Struct()
    output.update({"eventUrl": "repo_url", "event": "event"})
    
    logger.info(f"sending message {output}")

    await sdk.messaging.send_output(output)
    logger.info("github repo event processed")


def finalizer(sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[GITHUB-TRIGGER FINALIZER]")
    logger.info("finalizing process...")


async def init():
    runner = await Runner().initialize()
    await runner.trigger_runner().with_initializer(initializer).with_runner(
        basic_runner
    ).with_finalizer(finalizer).run()


if __name__ == "__main__":
    loop = asyncio.new_event_loop()
    loop.run_until_complete(init())
    loop.run_forever()
    loop.close()
