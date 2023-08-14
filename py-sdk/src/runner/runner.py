import sys
from dataclasses import dataclass

from exceptions import FailedLoadingConfigError, JetStreamConnectionError, NATSConnectionError
from exit.exit_runner import ExitRunner
from loguru import logger
from loguru._logger import Logger
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from task.task_runner import TaskRunner
from trigger.trigger_runner import TriggerRunner
from vyper import v


@dataclass
class Runner:
    nc: NatsClient = NatsClient()
    js: JetStreamContext = None
    logger: Logger = None

    def __post_init__(self):
        self.initialize_config()
        self.initialize_logger()

    async def initialize(self):
        try:
            await self.nc.connect(v.get("nats.url"))
        except Exception as e:
            self.logger.error(f"error connecting to nats: {e}")
            raise NATSConnectionError(e)

        try:
            self.js = await self.nc.jetstream()
        except Exception as e:
            self.logger.error(f"error connecting to jetstream: {e}")
            raise JetStreamConnectionError(e)

    def initialize_config(self):
        v.automatic_env()

        if v.is_set("APP_CONFIG_PATH"):
            v.add_config_path(v.get("APP_CONFIG_PATH"))

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
            v.add_config_path(v.get("APP_CONFIG_PATH"))

        try:
            v.merge_in_config()
        except Exception as e:
            error = e

        if len(v.all_keys()) == 0:
            logger.error("no configuration found")
            raise FailedLoadingConfigError(error)

        v.set_default("metadata.base_path", "/")
        v.set_default("runner.logger.level", "INFO")
        v.set_default("runner.logger.output_paths", ["stdout"])
        v.set_default("runner.logger.error_output_paths", ["stderr"])

    def initialize_logger(self):
        logger_format = "<green>{time}</green> <level>{extra[context]} {message}</level>"
        output_paths = v.get("runner.logger.output_paths")
        error_output_paths = v.get("runner.logger.error_output_paths")

        logger.remove()  # Remove the pre-configured handler
        for output_path in output_paths:
            if output_path == "stdout":
                output_path = sys.stdout

            logger.add(
                output_path,
                colorize=True,
                format=logger_format,
                backtrace=False,
                diagnose=False,
                level=v.get("runner.logger.level"),
            )
        for error_output_path in error_output_paths:
            if output_path == "stderr":
                output_path = sys.stderr

            logger.add(
                error_output_path,
                colorize=True,
                format=logger_format,
                backtrace=True,
                diagnose=True,
                level="ERROR",
            )

        self.logger = logger.bind(context="[RUNNER CONFIG]")
        self.logger.info("logger initialized")

    async def trigger_runner(self):
        return await TriggerRunner(self.js, self.logger).initialize()

    async def task_runner(self):
        return await TaskRunner(self.js, self.logger).initialize()

    async def exit_runner(self):
        return await ExitRunner(self.js, self.logger).initialize()
