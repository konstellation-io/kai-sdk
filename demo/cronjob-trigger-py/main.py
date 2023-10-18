import asyncio
import uuid

from google.protobuf.wrappers_pb2 import StringValue
from runner.runner import Runner
from runner.trigger import trigger_runner
from sdk import kai_sdk as sdk
from sdk.centralized_config.centralized_config import Scope


async def initializer(kaiSDK: sdk.KaiSDK):
    logger = kaiSDK.logger.bind(context="[CRONJOB INITIALIZER]")
    logger.info("starting example...")
    kaiSDK.logger.info(f"process {kaiSDK.get_process()}")
    kaiSDK.logger.info(f"product {kaiSDK.get_product()}")
    kaiSDK.logger.info(f"workflow {kaiSDK.get_workflow()}")
    kaiSDK.logger.info(f"global {kaiSDK.get_global()}")
    kaiSDK.logger.info(f"version {kaiSDK.get_version()}")
    kaiSDK.logger.info(
        f"kv_product {kaiSDK.get_product_centralized_configuration_name()}"
    )
    kaiSDK.logger.info(
        f"kv_workflow {kaiSDK.get_workflow_centralized_configuration_name()}"
    )
    kaiSDK.logger.info(
        f"kv_process {kaiSDK.get_process_centralized_configuration_name()}"
    )
    kaiSDK.logger.info(
        f"kv_global {kaiSDK.get_global_centralized_configuration_name()}"
    )
    kaiSDK.logger.info(f"object-store {kaiSDK.get_object_store_name()}")

    kaiSDK.path_utils.logger.info(f"base path {kaiSDK.path_utils.get_base_path()}")
    kaiSDK.path_utils.logger.info(
        f"compose base path {kaiSDK.path_utils.compose_path('test')}"
    )

    value1 = await kaiSDK.centralized_config.get_config("test1")
    value2 = await kaiSDK.centralized_config.get_config("test2")

    await kaiSDK.centralized_config.set_config("test", "value", Scope.WorkflowScope)

    await kaiSDK.object_store.save("test", bytes("value-obj", "utf-8"))

    kaiSDK.centralized_config.logger.info(
        f"config values from comfig.yaml test1: {value1} test2: {value2}"
    )


async def cronjob_runner(
    trigger_runner: trigger_runner.TriggerRunner, kaiSDK: sdk.KaiSDK
):
    while True:
        logger = kaiSDK.logger.bind(context="[CRONJOB RUNNER]")
        logger.info("executing example...")
        request_id = str(uuid.uuid4())
        await kaiSDK.messaging.send_output_with_request_id(
            StringValue(value="hello world"), request_id
        )
        logger.info(f"waiting response for request id {request_id}...")
        response = await trigger_runner.get_response_channel(request_id)
        logger.info(f"response: {response}")
        await asyncio.sleep(3)


def finalizer(kaiSDK: sdk.KaiSDK):
    logger = kaiSDK.logger.bind(context="[CRONJOB FINALIZER]")
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
