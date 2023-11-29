import asyncio

from google.protobuf.any_pb2 import Any
from google.protobuf.wrappers_pb2 import StringValue
from runner.runner import Runner
from sdk import kai_sdk as sdk
from sdk.metadata.metadata import Metadata
from sdk.predictions.types import Filter, TimestampRange, UpdatePayloadFunc, Payload
from datetime import datetime
from typing import Any


async def initializer(kai_sdk: sdk.KaiSDK):
    logger = kai_sdk.logger.bind(context="[REDIS INITIALIZER]")
    logger.info("starting example...")


async def handler(kai_sdk: sdk.KaiSDK, response: Any):
    logger = kai_sdk.logger.bind(context="[REDIS HANDLER]")
    string_value = StringValue()

    response.Unpack(string_value)
    message = string_value.value
    logger.info(f"received message {message}")

    logger.info("saving prediction")
    key = "test"
    prediction = {"test": "test"}
    kai_sdk.predictions.save(key, prediction)
    logger.info(f"prediction saved with key {key} and value {prediction}")
    response = kai_sdk.predictions.get(key)
    logger.info(f"prediction retrieved with key {key} and value {response}")

    def func(payload: UpdatePayloadFunc) -> Payload:
        payload["test"] = "updatedValueByFunc"
        return payload

    kai_sdk.predictions.update(key, func)
    logger.info(f"prediction updated with key {key}")
    response = kai_sdk.predictions.get(key)
    logger.info(f"prediction retrieved with key {key} and value {response}")

    filter_ = Filter(
        version=Metadata.get_version(),
        workflow=Metadata.get_workflow(),
        workflow_type=Metadata.get_workflow_type(),
        process=Metadata.get_process(),
        request_id=kai_sdk.get_request_id(),
        creation_date=TimestampRange(start_date=datetime.fromisoformat("2023-11-28").timestamp(), end_date=datetime.now().timestamp()),
    )
    response = kai_sdk.predictions.find(filter=filter_)
    logger.info(f"prediction retrieved with filter {filter_} and value {response}")

    composed_string = f"{message}, processed by the task process!"
    logger.info(f"sending message {composed_string}")
    await kai_sdk.messaging.send_output(StringValue(value=composed_string))


def preprocessor(kai_sdk: sdk.KaiSDK, response: Any):
    logger = kai_sdk.logger.bind(context="[REDIS PREPROCESSOR]")
    logger.info("I am a sync preprocessor")


async def postprocessor(kai_sdk: sdk.KaiSDK, response: Any):
    logger = kai_sdk.logger.bind(context="[REDIS POSTPROCESSOR]")
    logger.info("I am an async postprocessor")
    await asyncio.sleep(0.0000001)
    logger.info("I am an async postprocessor, after sleep")


def finalizer(kai_sdk: sdk.KaiSDK):
    logger = kai_sdk.logger.bind(context="[REDIS FINALIZER]")
    logger.info("finalizing example synchronously...")


async def init():
    runner = await Runner().initialize()
    await runner.task_runner().with_initializer(initializer).with_handler(
        handler
    ).with_finalizer(finalizer).with_preprocessor(preprocessor).with_postprocessor(
        postprocessor
    ).run()


if __name__ == "__main__":
    loop = asyncio.new_event_loop()
    loop.run_until_complete(init())
    loop.run_forever()
    loop.close()
