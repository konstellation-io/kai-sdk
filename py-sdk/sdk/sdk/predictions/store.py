from __future__ import annotations

import json
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from datetime import datetime

import loguru
from loguru import logger
from redis import Redis
from vyper import v

from sdk.metadata.metadata import Metadata
from sdk.predictions.exceptions import (
    FailedToFindPredictionsError,
    FailedToGetPredictionError,
    FailedToInitializePredictionsStoreError,
    FailedToParseResultError,
    FailedToSavePredictionError,
    NotFoundError,
)
from sdk.predictions.types import Filter, Prediction


@dataclass
class PredictionsABC(ABC):
    @abstractmethod
    def save(self, id: str, value: dict[str, str]) -> None:
        pass

    @abstractmethod
    def get(self, id: str) -> Prediction:
        pass

    @abstractmethod
    def find(self, filter: Filter) -> list[Prediction]:
        pass


@dataclass
class Predictions(PredictionsABC):
    logger: loguru.Logger = field(init=False)
    request_id: str = field(init=False)
    client: Redis = field(init=False)

    def __post_init__(self):
        origin = logger._core.extra["origin"]
        self.logger = logger.bind(context=f"{origin}.[PREDICTIONS STORE]")
        try:
            self.client = Redis(
                host=v.get_string("predictions.endpoint"),
                username=v.get_string("predictions.username"),
                password=v.get_string("predictions.password"),
            )
        except Exception as e:
            self.logger.error(f"failed to initialize predictions store: {e}")
            raise FailedToInitializePredictionsStoreError(e)

        self.logger.info("successfully initialized predictions store")

    def save(self, id: str, value: dict[str, str]) -> None:
        try:
            Prediction(
                timestamp=datetime.now().isoformat(),
                payload=value,
                metadata={
                    "product": Metadata.get_product(),
                    "version": Metadata.get_version(),
                    "workflow": Metadata.get_workflow(),
                    "process": Metadata.get_process(),
                    "request_id": self.request_id,
                },
            )
            self.client.hset(id, value)
        except Exception as e:
            self.logger.error(f"failed to save prediction with {id} to the predictions store: {e}")
            raise FailedToSavePredictionError(id, e)

        self.logger.info(f"successfully saved prediction with {id} to the predictions store")

    def get(self, id: str) -> Prediction:
        try:
            prediction = self.client.get(id)
        except Exception as e:
            self.logger.error(f"failed to get prediction {id} from the predictions store: {e}")
            raise FailedToGetPredictionError(id, e)

        self.logger.info(f"successfully got prediction {id} from the predictions store")

        if not prediction:
            self.logger.error(f"prediction {id} not found in the predictions store")
            raise NotFoundError(id)

        self.logger.info(f"successfully found prediction {id} from the predictions store")
        return self._parse_result(prediction)

    def find(self, filter: Filter) -> list[Prediction]:
        try:
            predictions = self.client.get(filter)
        except Exception as e:
            self.logger.error(
                f"failed to find predictions from the predictions store matching the filter {filter}: {e}"
            )
            raise FailedToFindPredictionsError(filter, e)

        self.logger.info(f"successfully found predictions from the predictions store matching the filter {filter}")
        return [self._parse_result(prediction) for prediction in predictions]

    def _parse_result(self, result: dict[str, str]) -> Prediction:
        try:
            result = json.loads(result)
            return Prediction(
                timestamp=result.get("timestamp", ""),
                result=result.get("result", {}),
                metadata=result.get("metadata", {}),
            )
        except Exception as e:
            self.logger.error(f"failed to parse result {result}: {e}")
            raise FailedToParseResultError(result, e)
