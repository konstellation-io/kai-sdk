from typing import Optional

from sdk.predictions.types import Filter


class FailedToInitializePredictionsStoreError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = "failed to initialize predictions store"
        super().__init__(f"{message}: {error}" if error else message)


class FailedToSavePredictionError(Exception):
    def __init__(self, key: str, error: Optional[Exception] = None):
        message = f"failed to save prediction with {key} to the predictions store"
        super().__init__(f"{message}: {error}" if error else message)


class FailedToGetPredictionError(Exception):
    def __init__(self, key: str, error: Optional[Exception] = None):
        message = f"failed to get prediction {key} from the predictions store"
        super().__init__(f"{message}: {error}" if error else message)


class NotFoundError(Exception):
    def __init__(self, key: str, error: Optional[Exception] = None):
        message = f"prediction {key} not found in the predictions store"
        super().__init__(f"{message}: {error}" if error else message)


class FailedToFindPredictionsError(Exception):
    def __init__(self, filter: Filter, error: Optional[Exception] = None):
        message = f"failed to find predictions from the predictions store matching the filter {filter}"
        super().__init__(f"{message}: {error}" if error else message)


class FailedToParseResultError(Exception):
    def __init__(self, result: dict[str, str], error: Optional[Exception] = None):
        message = f"failed to parse result {result}"
        super().__init__(f"{message}: {error}" if error else message)


class MissingRequiredFilterFieldError(Exception):
    def __init__(self, field: str, error: Optional[Exception] = None):
        message = f"filter {field} is required"
        super().__init__(f"{message}: {error}" if error else message)
