from __future__ import annotations

from abc import ABC
from dataclasses import dataclass

import loguru
from loguru import logger


@dataclass
class PersistentStorageABC(ABC):
    pass


@dataclass
class PersistentStorage(PersistentStorageABC):
    logger: loguru.Logger = logger.bind(context="[PERSISTENT STORAGE]")
