from typing import Optional


class FailedInitializingConfigError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = "failed initializing configuration"
        super().__init__(f"{message}: {error}" if error else message)


class FailedGettingConfigError(Exception):
    def __init__(self, key: str, scope: str, error: Optional[Exception] = None):
        message = f"failed getting configuration given key {key} and scope {scope}"
        super().__init__(f"{message}: {error}" if error else message)


class FailedSettingConfigError(Exception):
    def __init__(self, key: str, scope: str, error: Optional[Exception] = None):
        message = f"failed setting configuration given key {key} and scope {scope}"
        super().__init__(f"{message}: {error}" if error else message)


class FailedDeletingConfigError(Exception):
    def __init__(self, key: str, scope: str, error: Optional[Exception] = None):
        message = f"failed deleting configuration given key {key} and scope {scope}"
        super().__init__(f"{message}: {error}" if error else message)
