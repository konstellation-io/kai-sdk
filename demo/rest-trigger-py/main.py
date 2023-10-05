import asyncio
import signal
import uuid
from concurrent.futures import ThreadPoolExecutor

import uvicorn
from fastapi import Depends, FastAPI
from google.protobuf.wrappers_pb2 import StringValue
from runner.runner import Runner
from runner.trigger import trigger_runner
from sdk import kai_sdk as sdk
import sys
app = FastAPI()
import contextvars

sdk_instance = contextvars.ContextVar("sdk_instance")


async def initializer(sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[REST SERVER INITIALIZER]")
    logger.info("starting example...")
    sdk_instance.set(sdk)


async def rest_server_runner(
    trigger_runner: trigger_runner.TriggerRunner, sdk: sdk.KaiSDK
):
    logger = sdk.logger.bind(context="[REST SERVER RUNNER]")
    logger.info("executing example...")

    logger.info("starting rest server port 8080")

    executor = ThreadPoolExecutor(max_workers=1)
    future = executor.submit(init_server)

    def shutdown():
        print("Shutting down server...")
        executor.shutdown(wait=False)
        loop.stop()
        sys.exit(0)

    signal.signal(signal.SIGINT, lambda s, f: shutdown())
    signal.signal(signal.SIGTERM, lambda s, f: shutdown())

    future.result()  # This blocks until the server is stopped


def finalizer(sdk: sdk.KaiSDK):
    logger = sdk.logger.bind(context="[REST SERVER FINALIZER]")
    logger.info("finalizing example...")


async def init():
    runner = await Runner().initialize()
    await runner.trigger_runner().with_initializer(initializer).with_runner(
        rest_server_runner
    ).with_finalizer(finalizer).run()


async def response_handler(name: str) -> str:
    sdk_: sdk = sdk_instance.get()
    logger = sdk.logger.bind(context="[RESPONSE HANDLER]")
    logger.info(f"response handler received {name=}")

    response = StringValue(value=f"Hello {name}")
    request_id = str(uuid.uuid4())

    sdk_.logger.info(f"response handler returning {response=}")

    await sdk_.messaging.send_output_with_request_id(response, request_id)
    logger.info(f"waiting response for request id {request_id}...")

    response = await sdk_.trigger_runner.get_response_channel(request_id)
    logger.info(f"response: {response}")

    return response


def init_server():
    uvicorn.run(app, host="0.0.0.0", port=8080)


@app.get("/hello", response_model=dict)
async def hello(name: str = Depends(response_handler)):
    return {"message": name.split(",")}


if __name__ == "__main__":
    loop = asyncio.new_event_loop()
    loop.run_until_complete(init())
    loop.run_forever()
    loop.close()
