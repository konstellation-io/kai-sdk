class FailedGettingMaxMessageSizeError(Exception):
    def __init__(self, error: Exception = None):
        message = "failed getting max message size"
        super().__init__(f"{message}: {error}" if error else message)


class MessageTooLargeError(Exception):
    def __init__(self, message_size: int, max_message_size: int):
        super().__init__(f"message size {message_size} is larger than max message size {max_message_size}")
