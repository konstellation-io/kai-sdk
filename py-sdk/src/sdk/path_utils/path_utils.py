import os
from dataclasses import dataclass

from loguru import logger
from loguru._logger import Logger
from vyper import v


@dataclass
class PathUtils:
    logger: Logger = logger.bind(context="[PATH UTILS]")

    @staticmethod
    def get_base_path() -> str:
        return v.get("metadata.base_path")

    @staticmethod
    def compose_path(*relative_path)-> str:
        base_path = PathUtils.get_base_path()
        if not relative_path:
            return base_path

        return os.path.join(base_path, *relative_path)
