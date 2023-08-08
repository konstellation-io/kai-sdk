class UndefinedObjectStoreError(Exception):
    def __init__(self):
        super().__init__("undefined object store")


class FailedObjectStoreInitializationError(Exception):
    def __init__(self, error: Exception = None):
        message = "failed object store initialization"
        super().__init__(f"{message}: {error}" if error else message)


class EmptyPayloadError(Exception):
    def __init__(self):
        super().__init__("empty payload")


class FailedCompilingRegexpError(Exception):
    def __init__(self, error: Exception = None):
        message = "error compiling regexp"
        super().__init__(f"{message}: {error}" if error else message)


class FailedListingFilesError(Exception):
    def __init__(self, error: Exception = None):
        message = "failed listing objects from the object store"
        super().__init__(f"{message}: {error}" if error else message)


class FailedGettingFileError(Exception):
    def __init__(self, key: str, error: Exception = None):
        message = f"failed getting file {key} from the object store"
        super().__init__(f"{message}: {error}" if error else message)


class FailedSavingFileError(Exception):
    def __init__(self, key: str, error: Exception = None):
        message = f"failed saving file {key} to the object store"
        super().__init__(f"{message}: {error}" if error else message)


class FailedDeletingFileError(Exception):
    def __init__(self, key: str, error: Exception = None):
        message = f"failed deleting file {key} from the object store"
        super().__init__(f"{message}: {error}" if error else message)


class FailedPurgingFilesError(Exception):
    def __init__(self, error: Exception = None):
        message = "failed purging objects from the object store"
        super().__init__(f"{message}: {error}" if error else message)
