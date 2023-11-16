from __future__ import annotations

import sys
from dataclasses import dataclass, field
from functools import reduce
from datetime import datetime

import loguru
from loguru import logger
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from vyper import v

from runner.exceptions import FailedLoadingConfigError, JetStreamConnectionError, NATSConnectionError
from runner.exit.exit_runner import ExitRunner
from runner.task.task_runner import TaskRunner
from runner.trigger.trigger_runner import TriggerRunner
import json


def sink_serializer(message):
    record = message.record
    time = datetime.utcfromtimestamp(record["time"].timestamp()).isoformat(timespec="milliseconds") + "Z"
    filepath = record["file"].path
    filepath = filepath.split("py-sdk/sdk/")[1] if "py-sdk/sdk/" in filepath else filepath.split("py-sdk/runner/")[1]
    filepath = filepath + ":" + str(record["line"])
    simplified = {
        "L": record["level"].name,
        "T": time,
        "N": record["extra"]["context"],
        "C": filepath,
        "M": record["message"],
        "request_id": record["extra"]["request_id"],
    }
    serialized = json.dumps(simplified)
    print(serialized)

LOGGER_FORMAT = (
    "<green>{time:YYYY-MM-DDTHH:mm:ss.SSS}Z</green> "
    "<cyan>{level}</cyan> {extra[context]} <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> "
    "<level>{message}</level> <level>{extra[request_id]}</level>"
)

MANDATORY_CONFIG_KEYS = [
    "metadata.product_id",
    "metadata.workflow_name",
    "metadata.process_name",
    "metadata.version_tag",
    "metadata.base_path",
    "nats.url",
    "nats.stream",
    "nats.output",
    "centralized_configuration.global.bucket",
    "centralized_configuration.product.bucket",
    "centralized_configuration.workflow.bucket",
    "centralized_configuration.process.bucket",
    "minio.endpoint",
    "minio.client_user",  # generated user for the bucket
    "minio.client_password",  # generated user's password for the bucket
    "minio.ssl",  # Enable or disable SSL
    "minio.bucket",  # Bucket to be used
    "auth.endpoint",  # keycloak endpoint
    "auth.client",  # Client to be used to authenticate
    "auth.client_secret",  # Client's secret to be used
    "auth.realm",  # Realm
]


@dataclass
class Runner:
    nc: NatsClient = NatsClient()
    js: JetStreamContext = field(init=False)
    logger: loguru.Logger = field(init=False)

    def __post_init__(self) -> None:
        self.initialize_config()
        self.initialize_logger()

    async def initialize(self) -> Runner:
        try:
            self.js = self.nc.jetstream()
        except Exception as e:
            self.logger.error(f"error connecting to jetstream: {e}")
            raise JetStreamConnectionError(e)

        try:
            await self.nc.connect(v.get_string("nats.url"))
        except Exception as e:
            self.logger.error(f"error connecting to nats: {e}")
            raise NATSConnectionError(e)

        return self

    def _validate_config(self, keys: dict[str]) -> None:
        for key in MANDATORY_CONFIG_KEYS:
            try:
                _ = reduce(lambda d, k: d[k], key.split("."), keys)
            except Exception:
                raise FailedLoadingConfigError(Exception(f"missing mandatory configuration key: {key}"))

    def initialize_config(self) -> None:
        v.set_env_prefix("KAI")
        v.automatic_env()

        if v.is_set("APP_CONFIG_PATH"):
            v.add_config_path(v.get_string("APP_CONFIG_PATH"))

        v.set_config_name("config")
        v.set_config_type("yaml")
        v.add_config_path(".")

        error = None
        try:
            v.read_in_config()
        except Exception as e:
            error = e

        v.set_config_name("app")
        v.set_config_type("yaml")
        v.add_config_path(".")

        if v.is_set("APP_CONFIG_PATH"):
            v.add_config_path(v.get_string("APP_CONFIG_PATH"))

        try:
            v.merge_in_config()
        except Exception as e:
            error = e

        if len(v.all_keys()) == 0:
            raise FailedLoadingConfigError(error)

        self._validate_config(v.all_settings())

        v.set_default("metadata.base_path", "/")
        v.set_default("runner.subscriber.ack_wait_time", 22)
        v.set_default("runner.logger.level", "INFO")
        v.set_default("runner.logger.output_paths", ["stdout"])
        v.set_default("runner.logger.error_output_paths", ["stderr"])

    def initialize_logger(self) -> None:
        output_paths = v.get("runner.logger.output_paths")
        error_output_paths = v.get("runner.logger.error_output_paths")

        logger.remove()  # Remove the pre-configured handler
        for output_path in output_paths:
            if output_path == "stdout" or output_path == "console":
                output_path = sys.stdout
            elif output_path == "json":
                output_path = sink_serializer

            logger.add(
                output_path,
                colorize=True,
                format=LOGGER_FORMAT,
                backtrace=False,
                diagnose=False,
                level=v.get_string("runner.logger.level"),
            )

        for error_output_path in error_output_paths:
            if error_output_path == "stderr" or error_output_path == "console":
                error_output_path = sys.stderr
            elif error_output_path == "json":
                error_output_path = sink_serializer

            logger.add(
                error_output_path,
                colorize=True,
                format=LOGGER_FORMAT,
                backtrace=True,
                diagnose=True,
                level="ERROR",
            )

        logger.configure(extra={"context": "[UNKNOWN]", "request_id": "{}"})

        self.logger = logger.bind(context="[RUNNER CONFIG]")
        self.logger.info("logger initialized")

    def trigger_runner(self) -> TriggerRunner:
        return TriggerRunner(self.nc, self.js)

    def task_runner(self) -> TaskRunner:
        return TaskRunner(self.nc, self.js)

    def exit_runner(self) -> ExitRunner:
        return ExitRunner(self.nc, self.js)
