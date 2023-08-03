from messaging_utils import size_in_mb


class FailedGettingMaxMessageSizeError(Exception):
    def __init__(self, error=None):
        message = "failed getting max message size"
        super().__init__(f"{message}: {error}" if error else message)


class MessageTooLargeError(Exception):
    def __init__(self, message_size, max_message_size):
        super().__init__(
            f"message size {size_in_mb(message_size)} is larger than max message size {size_in_mb(max_message_size)}"
        )
