from typing import Optional


class UndefinedRunnerFunctionError(Exception):
    def __init__(self):
        message = "undefined runner function"
        super().__init__(message)


class NewRequestMsgError(Exception):
    def __init__(self, error: Optional[Exception] = None):
        message = "error creating new request message"
        super().__init__(f"{message}: {error}" if error else message)


class UndefinedResponseHandlerError(Exception):
    def __init__(self):
        message = "undefined response handler"
        super().__init__(message)


class HandlerError(Exception):
    def __init__(self, node_from: str, node_to: str, error: Optional[Exception] = None):
        message = f"error in node {node_from} executing handler for node {node_to}"
        super().__init__(f"{message}: {error}" if error else message)
