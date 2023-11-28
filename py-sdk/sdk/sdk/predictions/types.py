from dataclasses import dataclass


@dataclass
class Prediction:
    creation_date: float
    last_modified: float
    payload: dict[str, str]
    metadata: dict[str, str]


@dataclass
class TimestampRange:
    start_date: str
    end_date: str


@dataclass
class Filter:
    version: str
    workflow: str
    workflow_type: str
    process: str
    request_id: str
    timestamp: TimestampRange
