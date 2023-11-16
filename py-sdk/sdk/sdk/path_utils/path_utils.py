import os
from abc import ABC, abstractmethod
from dataclasses import dataclass

from vyper import v


@dataclass
class PathUtilsABC(ABC):
    @staticmethod
    @abstractmethod
    def get_base_path() -> str:
        pass

    @staticmethod
    @abstractmethod
    def compose_path(*relative_path: str) -> str:
        pass


@dataclass
class PathUtils(PathUtilsABC):
    @staticmethod
    def get_base_path() -> str:
        return v.get_string("metadata.base_path")

    @staticmethod
    def compose_path(*relative_path: str) -> str:
        base_path = PathUtils.get_base_path()
        return os.path.join(base_path, *relative_path) if relative_path else base_path
