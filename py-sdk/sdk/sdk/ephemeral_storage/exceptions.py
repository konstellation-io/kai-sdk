from typing import Optional


class UndefinedEphemeralStorageError(Exception):
    def __init__(self):
        super().__init__("undefined ephemeral storage")


class FailedEphemeralStorageInitializationError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = "failed ephemeral storage initialization"
        super().__init__(f"{message}: {error}" if error else message)


class FailedCompilingRegexpError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = "error compiling regexp"
        super().__init__(f"{message}: {error}" if error else message)


class FailedListingFilesError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = "failed listing objects from the ephemeral storage"
        super().__init__(f"{message}: {error}" if error else message)


class FailedGettingFileError(Exception):
    def __init__(self, key: str, error: Optional[Exception] = None):
        message = f"failed getting file {key} from the ephemeral storage"
        super().__init__(f"{message}: {error}" if error else message)


class FailedSavingFileError(Exception):
    def __init__(self, key: str, error: Optional[Exception] = None):
        message = f"failed saving file {key} to the ephemeral storage"
        super().__init__(f"{message}: {error}" if error else message)


class FailedDeletingFileError(Exception):
    def __init__(self, key: str, error: Optional[Exception] = None):
        message = f"failed deleting file {key} from the ephemeral storage"
        super().__init__(f"{message}: {error}" if error else message)


class FailedPurgingFilesError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = "failed purging objects from the ephemeral storage"
        super().__init__(f"{message}: {error}" if error else message)
