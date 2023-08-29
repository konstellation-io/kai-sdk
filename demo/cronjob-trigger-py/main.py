from sdk import kai_sdk as sdk
from runner.trigger import trigger_runner
from runner.runner import Runner
import asyncio


async def initializer(sdk: sdk.KaiSDK):
    sdk.metadata.logger.info(f"process {sdk.metadata.get_process()}")
    sdk.metadata.logger.info(f"product {sdk.metadata.get_product()}")
    sdk.metadata.logger.info(f"workflow {sdk.metadata.get_workflow()}")
    sdk.metadata.logger.info(f"version {sdk.metadata.get_version()}")
    sdk.metadata.logger.info(f"kv_product {sdk.metadata.get_key_value_store_product_name()}")
    sdk.metadata.logger.info(f"kv_workflow {sdk.metadata.get_key_value_store_workflow_name()}")
    sdk.metadata.logger.info(f"kv_process {sdk.metadata.get_key_value_store_process_name()}")
    sdk.metadata.logger.info(f"object-store {sdk.metadata.get_object_store_name()}")

    sdk.path_utils.logger.info(f"base path {sdk.path_utils.get_base_path()}")
    sdk.path_utils.logger.info(f"compose base path {sdk.path_utils.compose_path('test')}")

    value1 = await sdk.centralized_config.get_config("test1")
    value2 = await sdk.centralized_config.get_config("test2")

    sdk.centralized_config.logger.info(f"config values from comfig.yaml test1: {value1} test2: {value2}")


def cronjob_runner(trigger_runner: trigger_runner.TriggerRunner, sdk: sdk.KaiSDK):
    sdk.logger.info("Cronjob Runner")

def finalizer(ctx: sdk.KaiSDK):
    ctx.logger.info("Finalizer")


async def init():
    runner = await Runner().initialize()
    await runner.trigger_runner().with_initializer(initializer).with_runner(cronjob_runner).with_finalizer(finalizer).run()

if __name__ == "__main__":    
    loop = asyncio.get_event_loop()
    loop.run_until_complete(init())
    loop.run_forever()
    loop.close()