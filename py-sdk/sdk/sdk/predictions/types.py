from dataclasses import dataclass
from typing import Optional, Any, Callable

Payload = dict[str, Any]
UpdatePayloadFunc = Callable[[Payload], Payload]

@dataclass
class Prediction:
    creation_date: float
    last_modified: float
    payload: Payload
    metadata: dict[str, str]


@dataclass
class TimestampRange:
    start_date: float
    end_date: float


@dataclass
class Filter:
    creation_date: TimestampRange
    version: Optional[str] = None
    workflow: Optional[str] = None
    workflow_type: Optional[str] = None
    process: Optional[str] = None
    request_id: Optional[str] = None
