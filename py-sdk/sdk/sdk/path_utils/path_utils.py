from __future__ import annotations

import os
from dataclasses import dataclass

import loguru
from loguru import logger
from vyper import v


@dataclass
class PathUtils:
    logger: loguru.Logger = logger.bind(context="[PATH UTILS]")

    @staticmethod
    def get_base_path() -> str:
        return v.get("metadata.base_path")

    @staticmethod
    def compose_path(*relative_path: str) -> str:
        base_path = PathUtils.get_base_path()
        return os.path.join(base_path, *relative_path) if relative_path else base_path
