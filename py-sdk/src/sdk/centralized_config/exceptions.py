class FailedInitializingConfigError(Exception):
    def __init__(self, error: Exception):
        message = f"failed initializing configuration: {error}"
        super().__init__(message)


class FailedGettingConfigError(Exception):
    def __init__(self, key: str, scope: str, error: Exception):
        message = f"failed getting configuration given key {key} and scope {scope}: {error}"
        super().__init__(message)


class FailedSettingConfigError(Exception):
    def __init__(self, key: str, scope: str, error: Exception):
        message = f"failed setting configuration given key {key} and scope {scope}: {error}"
        super().__init__(message)


class FailedDeletingConfigError(Exception):
    def __init__(self, key: str, scope: str, error: Exception):
        message = f"failed deleting configuration given key {key} and scope {scope}: {error}"
        super().__init__(message)
