import asyncio
import uuid

from google.protobuf.wrappers_pb2 import StringValue
from runner.runner import Runner
from runner.trigger import trigger_runner
from sdk import kai_sdk as sdk
from sdk.centralized_config.centralized_config import Scope


async def initializer(kai_sdk: sdk.KaiSDK):
    logger = kai_sdk.logger.bind(context="[CRONJOB INITIALIZER]")
    logger.info("starting example...")
    kai_sdk.logger.info(f"process {kai_sdk.metadata.get_process()}")
    kai_sdk.logger.info(f"product {kai_sdk.metadata.get_product()}")
    kai_sdk.logger.info(f"workflow {kai_sdk.metadata.get_workflow()}")
    kai_sdk.logger.info(f"version {kai_sdk.metadata.get_version()}")
    kai_sdk.logger.info(
        f"kv_product {kai_sdk.metadata.get_product_centralized_configuration_name()}"
    )
    kai_sdk.logger.info(
        f"kv_workflow {kai_sdk.metadata.get_workflow_centralized_configuration_name()}"
    )
    kai_sdk.logger.info(
        f"kv_process {kai_sdk.metadata.get_process_centralized_configuration_name()}"
    )
    kai_sdk.logger.info(
        f"kv_global {kai_sdk.metadata.get_global_centralized_configuration_name()}"
    )
    kai_sdk.logger.info(f"object-store {kai_sdk.metadata.get_ephemeral_storage_name()}")

    kai_sdk.logger.info(f"base path {kai_sdk.path_utils.get_base_path()}")
    kai_sdk.logger.info(f"compose base path {kai_sdk.path_utils.compose_path('test')}")

    value1 = await kai_sdk.centralized_config.get_config("test1")
    value2 = await kai_sdk.centralized_config.get_config("test2")

    await kai_sdk.centralized_config.set_config("test", "value", Scope.WorkflowScope)

    await kai_sdk.storage.ephemeral.save("test", bytes("value-obj", "utf-8"))

    kai_sdk.logger.info(
        f"config values from config.yaml test1: {value1} test2: {value2}"
    )


async def cronjob_runner(
    trigger_runner: trigger_runner.TriggerRunner, kai_sdk: sdk.KaiSDK
):
    while True:
        logger = kai_sdk.logger.bind(context="[CRONJOB RUNNER]")
        logger.info("executing example...")
        request_id = str(uuid.uuid4())
        await kai_sdk.messaging.send_output_with_request_id(
            StringValue(value="hello world"), request_id
        )
        logger.info(f"waiting response for request id {request_id}...")
        response = await trigger_runner.get_response_channel(request_id)
        logger.info(f"response: {response}")
        await asyncio.sleep(3)


def finalizer(kai_sdk: sdk.KaiSDK):
    logger = kai_sdk.logger.bind(context="[CRONJOB FINALIZER]")
    logger.info("finalizing example...")


async def init():
    runner = await Runner().initialize()
    await runner.trigger_runner().with_initializer(initializer).with_runner(
        cronjob_runner
    ).with_finalizer(finalizer).run()


if __name__ == "__main__":
    loop = asyncio.new_event_loop()
    loop.run_until_complete(init())
    loop.run_forever()
    loop.close()
