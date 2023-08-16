class MissingRunnerFuncError(Exception):
    def __init__(self):
        message = "Runner function not found"
        super().__init__(message)
