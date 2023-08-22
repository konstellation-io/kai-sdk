class UndefinedDefaultHandlerFunctionError(Exception):
    def __init__(self):
        super().__init__("Undefined default handler")