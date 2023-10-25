from typing import Optional


class UndefinedPersistentStorageError(Exception):
    def __init__(self):
        super().__init__("undefined persistent storage")


class FailedPersistentStorageInitializationError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = "failed persistent storage initialization"
        super().__init__(f"{message}: {error}" if error else message)


class FailedListingFilesError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = "failed listing objects from the persistent storage"
        super().__init__(f"{message}: {error}" if error else message)


class FailedGettingFileError(Exception):
    def __init__(self, key: str, error: Optional[Exception] = None):
        message = f"failed getting file {key} from the persistent storage"
        super().__init__(f"{message}: {error}" if error else message)


class FailedSavingFileError(Exception):
    def __init__(self, key: str, error: Optional[Exception] = None):
        message = f"failed saving file {key} to the persistent storage"
        super().__init__(f"{message}: {error}" if error else message)


class FailedDeletingFileError(Exception):
    def __init__(self, key: str, error: Optional[Exception] = None):
        message = f"failed deleting file {key} from the persistent storage"
        super().__init__(f"{message}: {error}" if error else message)
