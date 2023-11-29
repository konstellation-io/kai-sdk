from typing import Optional


class EmptyNameError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = "the name cannot be empty"
        super().__init__(f"{message}: {error}" if error else message)


class InvalidVersionError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = "the given version is not valid, follow the semantic versioning specification"
        super().__init__(f"{message}: {error}" if error else message)


class EmptyModelError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = "the model cannot be empty"
        super().__init__(f"{message}: {error}" if error else message)


class FailedToSaveModelError(Exception):
    def __init__(self, name: str, error: Optional[Exception] = None):
        message = f"failed to save model {name} to the model registry"
        super().__init__(f"{message}: {error}" if error else message)


class FailedToGetModelError(Exception):
    def __init__(self, name: str, version: str, error: Optional[Exception] = None):
        message = f"failed to get model {name} with version {version} from the model registry"
        super().__init__(f"{message}: {error}" if error else message)


class FailedToListModelsError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = f"failed to list models from the model registry"
        super().__init__(f"{message}: {error}" if error else message)


class FailedToDeleteModelError(Exception):
    def __init__(self, name: str, error: Optional[Exception] = None):
        message = f"failed to delete model {name} from the model registry"
        super().__init__(f"{message}: {error}" if error else message)


class FailedToInitializeModelRegistryError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = "failed to initialize model registry"
        super().__init__(f"{message}: {error}" if error else message)


class MissingBucketError(Exception):
    def __init__(self, bucket: str):
        super().__init__(f"bucket: {bucket} not found in the persistent storage")


class FailedToListFilesError(Exception):
    def __init__(self, bucket: str, error: Optional[Exception] = None):
        message = f"failed to list files from the persistent storage bucket {bucket}"
        super().__init__(f"{message}: {error}" if error else message)


class FailedToGetFileError(Exception):
    def __init__(self, key: str, version: str, bucket: str, error: Optional[Exception] = None):
        message = f"failed to get file {key} with version {version} from the persistent storage bucket {bucket}"
        super().__init__(f"{message}: {error}" if error else message)


class FailedToSaveFileError(Exception):
    def __init__(self, key: str, bucket: str, error: Optional[Exception] = None):
        message = f"failed to save file {key} to the persistent storage bucket {bucket}"
        super().__init__(f"{message}: {error}" if error else message)


class FailedToDeleteFileError(Exception):
    def __init__(self, key: str, version: str, bucket: str, error: Optional[Exception] = None):
        message = f"failed to delete file {key} with version {version} from the persistent storage bucket {bucket}"
        super().__init__(f"{message}: {error}" if error else message)
