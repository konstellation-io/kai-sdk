import asyncio
import time

from google.protobuf.wrappers_pb2 import StringValue

from runner.runner import Runner
from runner.trigger import trigger_runner
from sdk import kai_sdk as sdk


async def initializer(sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[CRONJOB INITIALIZER]")
    logger.info("starting example...")
    sdk.metadata.logger.info(f"process {sdk.metadata.get_process()}")
    sdk.metadata.logger.info(f"product {sdk.metadata.get_product()}")
    sdk.metadata.logger.info(f"workflow {sdk.metadata.get_workflow()}")
    sdk.metadata.logger.info(f"version {sdk.metadata.get_version()}")
    sdk.metadata.logger.info(
        f"kv_product {sdk.metadata.get_key_value_store_product_name()}"
    )
    sdk.metadata.logger.info(
        f"kv_workflow {sdk.metadata.get_key_value_store_workflow_name()}"
    )
    sdk.metadata.logger.info(
        f"kv_process {sdk.metadata.get_key_value_store_process_name()}"
    )
    sdk.metadata.logger.info(f"object-store {sdk.metadata.get_object_store_name()}")

    sdk.path_utils.logger.info(f"base path {sdk.path_utils.get_base_path()}")
    sdk.path_utils.logger.info(
        f"compose base path {sdk.path_utils.compose_path('test')}"
    )

    value1 = await sdk.centralized_config.get_config("test1")
    value2 = await sdk.centralized_config.get_config("test2")

    sdk.centralized_config.logger.info(
        f"config values from comfig.yaml test1: {value1} test2: {value2}"
    )


async def cronjob_runner(trigger_runner: trigger_runner.TriggerRunner, sdk: sdk.KaiSDK):
    while True:
        logger = sdk.logger.bind(context="[CRONJOB RUNNER]")
        logger.info("executing example...")
        await sdk.messaging.send_output(StringValue(value="hello world"))
        time.sleep(2)


def finalizer(sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[CRONJOB FINALIZER]")
    logger.info("finalizing example...")


def handler(sdk: sdk.KaiSDK, response: StringValue):
    logger = sdk.logger.bind(context="[CRONJOB HANDLER]")
    logger.info(f"handling response: {response.value}")


async def init():
    runner = await Runner().initialize()
    await runner.trigger_runner().with_initializer(initializer).with_runner(
        cronjob_runner
    ).with_finalizer(finalizer).run()
    # await runner.task_runner().with_initializer(initializer).with_handler(handler).with_finalizer(finalizer).run()
    # await runner.exit_runner().with_initializer(initializer).with_handler(handler).with_finalizer(finalizer).run()


if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(init())
    loop.run_forever()
    loop.close()
