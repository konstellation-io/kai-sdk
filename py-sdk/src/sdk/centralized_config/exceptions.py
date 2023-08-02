class FailedGettingConfigGivenKey(Exception):
    def __init__(self, key: str):
        message = f"config not found in any key-value store for key {key}"
        super().__init__(message)

class FailedInitializingConfig(Exception):
    def __init__(self, error: Exception):
        message = f"failed initializing configuration: {error}"
        super().__init__(message)

class FailedGettingConfig(Exception):
    def __init__(self, key: str, scope: str, error: Exception):
        message = f"failed getting configuration given key {key} and scope {scope}: {error}"
        super().__init__(message)

class FailedSettingConfig(Exception):
    def __init__(self, key: str, scope: str, error: Exception):
        message = f"failed setting configuration given key {key} and scope {scope}: {error}"
        super().__init__(message)

class FailedDeletingConfig(Exception):
    def __init__(self, key: str, scope: str, error: Exception):
        message = f"failed deleting configuration given key {key} and scope {scope}: {error}"
        super().__init__(message)