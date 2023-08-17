class UndefinedRunnerFunctionError(Exception):
    def __init__(self):
        message = "Undefined runner function"
        super().__init__(message)
