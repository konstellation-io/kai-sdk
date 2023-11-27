from dataclasses import dataclass


@dataclass
class Prediction:
    timestamp: str
    payload: dict[str, str]
    metadata: dict[str, str]


@dataclass
class TimestampRange:
    start_date: str
    end_date: str


@dataclass
class Filter:
    request_id: str
    workflow: str
    process: str
    version: str
    timestamp: TimestampRange
