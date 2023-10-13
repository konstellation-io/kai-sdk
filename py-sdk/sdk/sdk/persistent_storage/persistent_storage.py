from __future__ import annotations

from abc import ABC
from dataclasses import dataclass


@dataclass
class PersistentStorageABC(ABC):
    pass


@dataclass
class PersistentStorage(PersistentStorageABC):
    pass