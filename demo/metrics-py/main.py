import asyncio

from google.protobuf.any_pb2 import Any
from google.protobuf.wrappers_pb2 import StringValue
from runner.runner import Runner
from sdk import kai_sdk as sdk

counter = None


async def initializer(kai_sdk: sdk.KaiSDK):
    logger = kai_sdk.logger.bind(context="[METRICS INITIALIZER]")
    logger.info("starting example...")
    logger.info("creating counter messages_received...")
    metrics_client = kai_sdk.measurements.get_metrics_client()
    global counter
    counter = metrics_client.create_counter(
        name="messages_received", description="messages received", unit="integer"
    )


async def handler(kai_sdk: sdk.KaiSDK, response: Any):
    logger = kai_sdk.logger.bind(context="[METRICS HANDLER]")

    string_value = StringValue()
    response.Unpack(string_value)
    message = string_value.value
    logger.info(f"received message {message} and incrementing counter by 1...")
    global counter
    counter.add(1)

    composed_string = f"{message}, processed by the metrics process!"
    logger.info(f"sending message {composed_string}")
    await kai_sdk.messaging.send_output(StringValue(value=composed_string))


def finalizer(kai_sdk: sdk.KaiSDK):
    logger = kai_sdk.logger.bind(context="[METRICS FINALIZER]")
    logger.info("finalizing example...")


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
