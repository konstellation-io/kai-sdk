from typing import Optional


class FailedGettingMaxMessageSizeError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = "failed getting max message size"
        super().__init__(f"{message}: {error}" if error else message)


class MessageTooLargeError(Exception):
    def __init__(self, message_size: str, max_message_size: str):
        super().__init__(f"message size {message_size} is larger than max message size {max_message_size}")
