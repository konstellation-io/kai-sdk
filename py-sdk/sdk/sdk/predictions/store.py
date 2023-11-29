from __future__ import annotations

import json
from abc import ABC, abstractmethod
from dataclasses import asdict, dataclass, field
from datetime import datetime

import loguru
from loguru import logger
from redis import Redis
from redis.commands.json.path import Path
from vyper import v

from sdk.metadata.metadata import Metadata
from sdk.predictions.exceptions import (
    FailedToFindPredictionsError,
    FailedToGetPredictionError,
    FailedToInitializePredictionsStoreError,
    FailedToParseResultError,
    FailedToSavePredictionError,
    MalformedEndpointError,
    MissingRequiredFilterFieldError,
    NotFoundError,
)
from sdk.predictions.types import Filter, Prediction


@dataclass
class PredictionsABC(ABC):
    @abstractmethod
    def save(self, id: str, function: callable) -> None:
        pass

    @abstractmethod
    def get(self, id: str) -> Prediction:
        pass

    @abstractmethod
    def find(self, filter: Filter) -> list[Prediction]:
        pass

    @abstractmethod
    def update(self, id: str, value: dict[str, str]) -> None:
        pass


@dataclass
class Predictions(PredictionsABC):
    logger: loguru.Logger = field(init=False)
    request_id: str = ""
    client: Redis = field(init=False)

    def __post_init__(self):
        origin = logger._core.extra["origin"]
        self.logger = logger.bind(context=f"{origin}.[PREDICTIONS STORE]")
        try:
            try:
                endpoint = v.get_string("predictions.endpoint")
                endpoint_ = endpoint.split(":")
                host = endpoint_[0]
                port = int(endpoint_[1])
            except Exception as e:
                self.logger.error(f"malformed endpoint: {e}")
                raise MalformedEndpointError(v.get_string("predictions.endpoint"), e)

            self.client = Redis(
                host=host,
                port=port,
                username=v.get_string("predictions.username"),
                password=v.get_string("predictions.password"),
            )
        except Exception as e:
            self.logger.error(f"failed to initialize predictions store: {e}")
            raise FailedToInitializePredictionsStoreError(e)

        self.logger.info("successfully initialized predictions store")

    def save(self, id: str, value: dict[str, str]) -> None:
        try:
            creation_timestamp = datetime.now().timestamp() * 1000  # milliseconds
            key = self._get_key_with_product_prefix(id)
            prediction = Prediction(
                creation_date=creation_timestamp,
                last_modified=creation_timestamp,
                payload=value,
                metadata={
                    "product": Metadata.get_product(),
                    "version": Metadata.get_version(),
                    "workflow": Metadata.get_workflow(),
                    "workflow_type": Metadata.get_workflow_type(),
                    "process": Metadata.get_process(),
                    "request_id": self.request_id,
                },
            )
            self.client.json().set(name=key, path=Path.root_path(), obj=asdict(prediction))
        except Exception as e:
            self.logger.error(f"failed to save prediction with {id} to the predictions store: {e}")
            raise FailedToSavePredictionError(id, e)

        self.logger.info(f"successfully saved prediction with {id} to the predictions store")

    def get(self, id: str) -> Prediction:
        try:
            key = self._get_key_with_product_prefix(id)
            prediction = self.client.json().get(key)
        except Exception as e:
            self.logger.error(f"failed to get prediction {id} from the predictions store: {e}")
            raise FailedToGetPredictionError(id, e)

        if not prediction:
            self.logger.error(f"prediction {id} not found in the predictions store")
            raise NotFoundError(id)

        self.logger.info(f"successfully found prediction {id} from the predictions store")
        return self._parse_result(prediction)

    def find(self, filter: Filter) -> list[Prediction]:
        self._validate_filter(filter)
        index = v.get_string("predictions.index_key")
        try:
            predictions = self.client.ft(index).search(query=self._build_query(filter))
        except Exception as e:
            self.logger.error(
                f"failed to find predictions from the predictions store matching the filter {filter}: {e}"
            )
            raise FailedToFindPredictionsError(filter, e)

        self.logger.info(f"successfully found predictions from the predictions store matching the filter {filter}")
        return [self._parse_result(prediction) for prediction in predictions]

    def update(self, id: str, function: callable) -> None:
        try:
            key = self._get_key_with_product_prefix(id)
            prediction = self.client.json().get(key)

            payload = prediction["payload"]
            new_payload = function(payload)
            last_modified = datetime.now().timestamp() * 1000  # milliseconds

            updated_prediction = Prediction(
                creation_date=prediction["creation_date"],
                last_modified=last_modified,
                payload=new_payload,
                metadata=prediction["metadata"],
            )

            self.client.json().set(name=key, path=Path.root_path(), obj=asdict(updated_prediction))
        except Exception as e:
            self.logger.error(f"failed to update prediction with {id} to the predictions store: {e}")
            raise FailedToSavePredictionError(id, e)

        self.logger.info(f"successfully updated prediction with {id} to the predictions store")

    def _parse_result(self, result: dict[str, str]) -> Prediction:
        try:
            return Prediction(**result)
        except Exception as e:
            self.logger.error(f"failed to parse result {result}: {e}")
            raise FailedToParseResultError(result, e)

    def _validate_filter(self, filter: Filter) -> None:
        if not filter.version:
            filter.version = Metadata.get_version()

        if not filter.timestamp:
            self.logger.error("filter timestamp is required")
            raise MissingRequiredFilterFieldError("timestamp")

        if not filter.timestamp.start_date:
            self.logger.error("filter timestamp start_date is required")
            raise MissingRequiredFilterFieldError("start_date")

        if not filter.timestamp.end_date:
            self.logger.error("filter timestamp end_date is required")
            raise MissingRequiredFilterFieldError("end_date")

    @staticmethod
    def _build_query(filter: Filter) -> str:
        query = f"@product:{Metadata.get_product()} @timestamp:[0 inf]"

        if filter.version:
            query = f"{query} @version:{filter.version}"

        if filter.workflow:
            query = f"{query} @workflow:{filter.workflow}"

        if filter.workflow_type:
            query = f"{query} @workflow_type:{filter.workflow_type}"

        if filter.process:
            query = f"{query} @process:{filter.process}"

        if filter.request_id:
            query = f"{query} @request_id:{filter.request_id}"

        if filter.timestamp:
            query = f"{query} @timestamp:[{filter.timestamp.start_date} {filter.timestamp.end_date}]"

        query = query.replace("-", "\\-")

        return query

    @staticmethod
    def _get_key_with_product_prefix(key: str) -> str:
        return f"{Metadata.get_product()}:{key}"
